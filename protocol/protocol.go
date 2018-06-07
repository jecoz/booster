/*
Copyright (C) 2018 Daniel Morandini

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package protocol provides functionalities to create payloads and compose packet
// that conform to the booster protocol.
package protocol

import (
	"fmt"
)

// Booster protocol version
const Version = "v0.2"

// Possible encodings
const (
	EncodingProtobuf uint8 = iota
	EncodingJson
)

const (
	CompressionNone uint8 = iota
)

const (
	EncryptionNone uint8 = iota
)

type Module string

// Module Identifiers
const (
	ModuleHeader  Module = "HE"
	ModulePayload Module = "PA"
)

// Tags used in the encoding and decoding of packets.
const (
	PacketOpeningTag = ">"
	PacketClosingTag = "<"
	ModuleOpeningTag = "["
	ModuleClosingTag = "]"
	Separator        = ":"
)

type Message int32

// Booster possible packet messages
const (
	// Commands
	MessageHello Message = iota
	MessageConnect
	MessageDisconnect
	MessageHeartbeat
	MessageNotify
	MessageMonitor
	MessageCtrl

	// Monitor related
	MessageNetworkStatus
	MessageNodeStatus
	MessageProxyUpdate
)

// Booster possible monitor features
const (
	MonitorProxy Message = iota
	MonitorNet
)

type Operation int32

// Tunnel operations
const (
	TunnelAdd Operation = iota
	TunnelAck
	TunnelRemove
)

// Ctrl operations
const (
	CtrlStop Operation = iota
	CtrlRestart
)

// OperationFromString converts raw, if possible, into a protocol
// known operation. Returns an error if no match is found.
func OperationFromString(raw string) (Operation, error) {
	switch raw {
	case "stop":
		return CtrlStop, nil
	case "restart":
		return CtrlRestart, nil
	default:
		return Operation(-1), fmt.Errorf("protocol: undefiled operation: %v", raw)
	}
}

// IsVersionSupported returns true if the current protocol version is compatible
// with the requested version.
func IsVersionSupported(v string) bool {
	// TODO(daniel): implement this check
	return v == Version
}
