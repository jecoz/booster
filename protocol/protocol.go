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
const Version = "v0.1.0"

// Possible encodings
const (
	EncodingProtobuf uint8 = 1
)

// Module Identifiers
const (
	ModuleHeader  string = "HE"
	ModulePayload        = "PA"
)

// Tags used in the encoding and decoding of packets.
const (
	PacketOpeningTag  = ">"
	PacketClosingTag  = "<"
	PayloadOpeningTag = "["
	PayloadClosingTag = "]"
	Separator         = ":"
)

type Message int32

// Booster possible packet messages
const (
	MessageHello      Message = 1
	MessageConnect            = 2
	MessageDisconnect         = 3
	MessageNode               = 4
	MessageHeartbeat          = 5
	MessageTunnel             = 6
	MessageNotify             = 7
	MessageInspect            = 8
	MessageBandwidth          = 9
	MessageCtrl               = 10
)

// Tunnel operations
const (
	TunnelAck    int32 = 1
	TunnelRemove       = 2
)

type Operation int32

// Ctrl operations
const (
	CtrlStop     Operation = 1
	CtrlRestart            = 2
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
