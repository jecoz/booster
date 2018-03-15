// Package provides the higher interface for dealing with booster instances
// that follow the booster protocol. It wraps togheter node, proxy, network.
package booster

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/danielmorandini/booster/network"
	"github.com/danielmorandini/booster/network/packet"
	"github.com/danielmorandini/booster/node"
	"github.com/danielmorandini/booster/protocol"
	"github.com/danielmorandini/booster/pubsub"
	"github.com/danielmorandini/booster/socks5"
)

// Proxy wraps that booster requires a proxy to implement.
type Proxy interface {
	// Notify acts as a one time registration function. The returned channel
	// should produce messages that give intormation about the state of the tunnels
	// generated by the proxy.
	Notify() (chan interface{}, error)

	// StopNotifying asks the receiver to stop generating tunnel notifications and
	// close the channel.
	StopNotifying(c chan interface{})

	// ListenAndServe should starts the actual proxy server, announcing it to the local
	// address.
	ListenAndServe(ctx context.Context, port int) error

	// Proto returns the string representation of the protocol used by the proxy.
	// Example: socks5.
	Proto() string
}

// PubSub describes the required functionalities of a publication/subscription object.
type PubSub interface {
	Sub(topic string) (chan interface{}, error)
	Unsub(c chan interface{}, topic string) error
	Pub(message interface{}, topic string) error
}

type SendConsumeCloser interface {
	SendCloser
	Consume() (<-chan *packet.Packet, error)
}

type SendCloser interface {
	Close() error
	Send(p *packet.Packet) error
}

// Booster wraps the parts that compose a booster node together.
type Booster struct {
	ID string

	*log.Logger
	Proxy Proxy
	PubSub

	Netconfig    network.Config
	stop         chan struct{}
	HeartbeatTTL time.Duration
}

// New creates a new configured booster node. Creates a network configuration
// based in the information contained in the protocol package.
//
// The internal proxy is configured to use the node dispatcher as network
// dialer.
func New(pport, bport int) (*Booster, error) {
	b := new(Booster)

	pp := strconv.Itoa(pport)
	bp := strconv.Itoa(bport)
	rn, err := node.New("localhost", pp, bp, true)
	if err != nil {
		return nil, err
	}

	id := sha1Hash([]byte(strconv.Itoa(pport)), []byte(strconv.Itoa(bport)))
	n := NewNet(rn, id)
	log := log.New(os.Stdout, "BOOSTER  ", log.LstdFlags)
	pubsub := pubsub.New()
	dialer := node.NewDispatcher(n)
	proxy := socks5.New(dialer)
	netconfig := network.Config{
		TagSet: packet.TagSet{
			PacketOpeningTag:  protocol.PacketOpeningTag,
			PacketClosingTag:  protocol.PacketClosingTag,
			PayloadClosingTag: protocol.PayloadClosingTag,
			Separator:         protocol.Separator,
		},
	}

	Nets.Set(id, n)

	b.ID = id
	b.Logger = log
	b.Proxy = proxy
	b.PubSub = pubsub
	b.Netconfig = netconfig
	b.stop = make(chan struct{})
	b.HeartbeatTTL = time.Second * 4

	return b, nil
}

// Run starts the proxy and booster node.
//
// This is a blocking routine that can be stopped using the Close() method.
// Traps INTERRUPT signals.
func (b *Booster) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errc := make(chan error, 3)
	_, pport, _ := net.SplitHostPort(Nets.Get(b.ID).LocalNode.PAddr.String())
	_, bport, _ := net.SplitHostPort(Nets.Get(b.ID).LocalNode.BAddr.String())
	pp, _ := strconv.Atoi(pport)
	bp, _ := strconv.Atoi(bport)
	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		errc <- b.ListenAndServe(ctx, bp)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		errc <- b.Proxy.ListenAndServe(ctx, pp)
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		errc <- b.UpdateRoot(ctx)
		wg.Done()
	}()

	// trap exit signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		for sig := range c {
			b.Printf("booster: signal (%v) received: exiting...", sig)
			b.Close()
			return
		}
	}()

	select {
	case err := <-errc:
		cancel()
		wg.Wait()
		return err
	case <-b.stop:
		cancel()
		wg.Wait()
		return fmt.Errorf("booster: stopped")
	}
}

// Close stops the Run routine. It drops the whole booster network, preparing for the
// node to reset or stop.
func (b *Booster) Close() error {
	b.stop <- struct{}{}
	return nil
}

// ListenAndServe shows to the network, listening for incoming tcp connections an
// turning them into booster connections.
func (b *Booster) ListenAndServe(ctx context.Context, port int) error {
	p := strconv.Itoa(port)
	ln, err := network.Listen("tcp", ":"+p, b.Netconfig)
	if err != nil {
		return err
	}
	defer ln.Close()

	b.Printf("listening on port: %v", p)

	errc := make(chan error)
	defer close(errc)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				errc <- fmt.Errorf("booster: cannot accept conn: %v", err)
				return
			}

			// send hello message first.
			if err := b.SendHello(ctx, conn); err != nil {
				errc <- err
				return
			}

			go b.Handle(ctx, conn)
		}
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		ln.Close()
		<-errc // wait for listener to return
		return ctx.Err()
	}
}

func (b *Booster) DialContext(ctx context.Context, netwrk, addr string) (*Conn, error) {
	dialer := network.NewDialer(new(net.Dialer), b.Netconfig)
	conn, err := dialer.DialContext(ctx, netwrk, addr)
	if err != nil {
		return nil, err
	}

	return b.RecvHello(ctx, conn)
}

func (b *Booster) Wire(ctx context.Context, network, target string) (*Conn, error) {
	// connect to the target node
	conn, err := b.DialContext(ctx, network, target)
	if err != nil {
		return nil, err
	}

	fail := func(err error) (*Conn, error) {
		conn.Close()
		return nil, err
	}

	// add the connection to the booster network. In case of an error,
	// i.e. such connection is already present, just close the underlying
	// network.Conn and return. Calling conn.Close() will close even the
	// other connection!
	err = Nets.Get(b.ID).AddConn(conn)
	if err != nil {
		conn.Conn.Close()
		return nil, err
	}

	// compose the notify packet which tells the receiver to start sending
	// information notifications when its state changes
	p := packet.New()
	enc := protocol.EncodingProtobuf
	h, err := protocol.TunnelNotifyHeader()
	_, err = p.AddModule(protocol.ModuleHeader, h, enc)
	if err != nil {
		return fail(err)
	}
	if err = conn.Send(p); err != nil {
		return fail(err)
	}

	b.Printf("booster: -> wire: %v", target)

	// inject the heartbeat message in the connection
	p, err = b.composeHeartbeat(nil)
	if err != nil {
		return fail(err)
	}
	if err = conn.Send(p); err != nil {
		return fail(err)
	}

	// start the timer that, when done, will close the connection if
	// no heartbeat message is received in time
	conn.HeartbeatTimer = time.AfterFunc(b.HeartbeatTTL, func() {
		b.Printf("booster: no heartbeat received from conn %v: timer expired", conn.ID)
		conn.Close()
	})

	// handle the newly added connection in a different goroutine.
	go b.Handle(ctx, conn)

	// set the connection as active
	conn.RemoteNode.SetIsActive(true)

	return conn, nil
}

func (b *Booster) UpdateRoot(ctx context.Context) error {
	errc := make(chan error)
	c, err := b.Proxy.Notify()
	if err != nil {
		return err
	}

	go func() {
		for i := range c {
			tm, ok := i.(socks5.TunnelMessage)
			if !ok {
				errc <- fmt.Errorf("unable to recognise workload message: %v", tm)
				return
			}
			node := Nets.Get(b.ID).LocalNode
			if err := b.UpdateNode(node, &tm, true); err != nil {
				b.Printf("booster: %v", err)
			}
		}

		errc <- nil
	}()

	defer close(errc)

	select {
	case err := <-errc:
		b.Proxy.StopNotifying(c)
		return err
	case <-ctx.Done():
		b.Proxy.StopNotifying(c)
		<-errc
		return ctx.Err()
	}
}

func sha1Hash(images ...[]byte) string {
	h := sha1.New()
	for _, image := range images {
		h.Write(image)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
