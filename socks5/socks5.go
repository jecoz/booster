// Package socks5 provides a SOCKS5 server implementation. See RFC 1928
// for protocol specification.
package socks5

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	socks5Version = uint8(5)
)

// Possible METHOD field values
const (
	socks5MethodNoAuth              = uint8(0)
	socks5MethodGSSAPI              = uint8(1)
	socks5MethodUsernamePassword    = uint8(2)
	socks5MethodNoAcceptableMethods = uint8(0xff)
)

// Possible CMD field values
const (
	socks5CmdConnect   = uint8(1)
	socks5CmdBind      = uint8(2)
	socks5CmdAssociate = uint8(3)
)

// Possible REP field values
const (
	socks5RespSuccess                 = uint8(0)
	socks5RespGeneralServerFailure    = uint8(1)
	socks5RespConnectionNotAllowed    = uint8(2)
	socks5RespNetworkUnreachable      = uint8(3)
	socks5RespHostUnreachable         = uint8(4)
	socks5RespConnectionRefused       = uint8(5)
	socks5RespTTLExpired              = uint8(6)
	socks5RespCommandNotSupported     = uint8(7)
	socks5RespAddressTypeNotSupported = uint8(8)
	socks5RespUnassigned              = uint8(9)
)

const (
	// FieldReserved should be used to fill fields marked as reserved.
	socks5FieldReserved = uint8(0x00)
)

const (
	// AddrTypeIPV4 is a version-4 IP address, with a length of 4 octets
	socks5IP4 = uint8(1)

	// AddrTypeFQDN field contains a fully-qualified domain name. The first
	// octet of the address field contains the number of octets of name that
	// follow, there is no terminating NUL octet.
	socks5FQDN = uint8(3)

	// AddrTypeIPV6 is a version-6 IP address, with a length of 16 octets.
	socks5IP6 = uint8(4)
)

var (
	supportedMethods = []uint8{socks5MethodNoAuth}
)

// Conn is a wrapper around io.ReadWriteCloser.
type Conn interface {
	io.ReadWriteCloser
}

// Dialer is the interface that wraps the DialContext function.
type Dialer interface {
	// DialContext opens a connection to addr, which should
	// be a canonical address with host and port.
	DialContext(ctx context.Context, network, addr string) (c net.Conn, err error)
}

// Socks5 represents a SOCKS5 proxy server implementation.
type Socks5 struct {
	*log.Logger

	// Dialer is used when connecting to a remote host. Could
	// be useful when chaining multiple proxies.
	Dialer
}

// ListenAndServe accepts and handles TCP connections
// using the SOCKS5 protocol.
func (s *Socks5) ListenAndServe(port int) error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	s.Printf("[TCP] listening on port: %v", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.Printf("[TCP Accept Error]: %v\n", err)
			continue
		}

		go func() {
			if err := s.Handle(conn); err != nil {
				s.Println(err)
			}
			s.Printf("[TCP]: connection to %v closed.\n", conn.RemoteAddr().String())
		}()
	}
}

// Handle performs the steps required to be SOCKS5 compliant.
// See RFC 1928 for details.
//
// Should run in its own go routine, closes the connection
// when returning.
func (s *Socks5) Handle(conn Conn) error {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	defer conn.Close()

	// method sub-negotiation phase
	if err := s.Negotiate(conn); err != nil {
		return err
	}

	// request details

	// len is jsut an estimation
	buf := make([]byte, 6+net.IPv4len)

	if _, err := io.ReadFull(conn, buf[:3]); err != nil {
		return errors.New("proxy: unable to read request: " + err.Error())
	}

	v := buf[0]   // protocol version
	cmd := buf[1] // command to execute
	_ = buf[2]    // reserved field

	// Check version number
	if v != socks5Version {
		return errors.New("proxy: unsupported version: " + string(v))
	}

	target, err := ReadAddress(conn)
	if err != nil {
		return err
	}

	var tconn Conn
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	switch cmd {
	case socks5CmdConnect:
		tconn, err = s.Connect(ctx, conn, target)
	case socks5CmdAssociate:
		tconn, err = s.Associate(ctx, conn, target)
	case socks5CmdBind:
		tconn, err = s.Bind(ctx, conn, target)
	default:
		return errors.New("unexpected CMD(" + strconv.Itoa(int(cmd)) + ")")
	}
	if err != nil {
		return err
	}
	defer tconn.Close()

	// start proxying
	go io.Copy(tconn, conn)
	io.Copy(conn, tconn)

	return nil
}

// ReadAddress reads hostname and port and converts them
// into its string format, properly formatted.
//
// r expects to read one byte that specifies the address
// format (1/3/4), follwed by the address itself and a
// 16 bit port number.
//
// addr == "" only when err != nil.
func ReadAddress(r io.Reader) (addr string, err error) {

	// cap is just an estimantion
	buf := make([]byte, 0, 2+net.IPv6len)
	buf = buf[:1]

	if _, err := io.ReadFull(r, buf); err != nil {
		return "", errors.New("proxy: unable to read address type: " + err.Error())
	}

	atype := buf[0] // address type

	bytesToRead := 0
	switch atype {
	case socks5IP4:
		bytesToRead = net.IPv4len
	case socks5IP6:
		bytesToRead = net.IPv6len
	case socks5FQDN:
		_, err := io.ReadFull(r, buf[:1])
		if err != nil {
			return "", errors.New("proxy: failed to read domain length: " + err.Error())
		}
		bytesToRead = int(buf[0])
	default:
		return "", errors.New("proxy: got unknown address type " + strconv.Itoa(int(atype)))
	}

	if cap(buf) < bytesToRead {
		buf = make([]byte, bytesToRead)
	} else {
		buf = buf[:bytesToRead]
	}
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", errors.New("proxy: failed to read address: " + err.Error())
	}

	var host string
	if atype == socks5FQDN {
		host = string(buf)
	} else {
		host = net.IP(buf).String()
	}

	if _, err := io.ReadFull(r, buf[:2]); err != nil {
		return "", errors.New("proxy: failed to read port: " + err.Error())
	}

	port := int(buf[0])<<8 | int(buf[1])
	addr = net.JoinHostPort(host, strconv.Itoa(port))

	return addr, nil
}

func (s *Socks5) getDialer() Dialer {
	if s.Dialer == nil {
		s.Dialer = new(net.Dialer)
	}

	return s.Dialer
}
