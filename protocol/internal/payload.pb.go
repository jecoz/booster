// Code generated by protoc-gen-go. DO NOT EDIT.
// source: payload.proto

package internal

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type PayloadCtrl struct {
	// operation is the operation that has to be performed.
	// See Ctrl operations in protocol.go
	Operation int32 `protobuf:"varint,1,opt,name=operation" json:"operation,omitempty"`
}

func (m *PayloadCtrl) Reset()                    { *m = PayloadCtrl{} }
func (m *PayloadCtrl) String() string            { return proto.CompactTextString(m) }
func (*PayloadCtrl) ProtoMessage()               {}
func (*PayloadCtrl) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *PayloadCtrl) GetOperation() int32 {
	if m != nil {
		return m.Operation
	}
	return 0
}

type PayloadBandwidth struct {
	// tot is the total number of bytes transmitted.
	Tot int64 `protobuf:"varint,1,opt,name=tot" json:"tot,omitempty"`
	// bandwidth is the current bandwidth.
	Bandwidth int64 `protobuf:"varint,2,opt,name=bandwidth" json:"bandwidth,omitempty"`
	// type is the transmission direction, i.e. dowload/upload
	Type string `protobuf:"bytes,3,opt,name=type" json:"type,omitempty"`
}

func (m *PayloadBandwidth) Reset()                    { *m = PayloadBandwidth{} }
func (m *PayloadBandwidth) String() string            { return proto.CompactTextString(m) }
func (*PayloadBandwidth) ProtoMessage()               {}
func (*PayloadBandwidth) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *PayloadBandwidth) GetTot() int64 {
	if m != nil {
		return m.Tot
	}
	return 0
}

func (m *PayloadBandwidth) GetBandwidth() int64 {
	if m != nil {
		return m.Bandwidth
	}
	return 0
}

func (m *PayloadBandwidth) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type PayloadInspect struct {
	// features contains the features that should be inspected.
	Features []int32 `protobuf:"varint,1,rep,packed,name=features" json:"features,omitempty"`
}

func (m *PayloadInspect) Reset()                    { *m = PayloadInspect{} }
func (m *PayloadInspect) String() string            { return proto.CompactTextString(m) }
func (*PayloadInspect) ProtoMessage()               {}
func (*PayloadInspect) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *PayloadInspect) GetFeatures() []int32 {
	if m != nil {
		return m.Features
	}
	return nil
}

type PayloadHello struct {
	// bport is the booster listening port.
	Bport string `protobuf:"bytes,1,opt,name=bport" json:"bport,omitempty"`
	// pport is the proxy listening port.
	Pport string `protobuf:"bytes,2,opt,name=pport" json:"pport,omitempty"`
}

func (m *PayloadHello) Reset()                    { *m = PayloadHello{} }
func (m *PayloadHello) String() string            { return proto.CompactTextString(m) }
func (*PayloadHello) ProtoMessage()               {}
func (*PayloadHello) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *PayloadHello) GetBport() string {
	if m != nil {
		return m.Bport
	}
	return ""
}

func (m *PayloadHello) GetPport() string {
	if m != nil {
		return m.Pport
	}
	return ""
}

type PayloadConnect struct {
	// target of the connect procedure.
	Target string `protobuf:"bytes,1,opt,name=target" json:"target,omitempty"`
}

func (m *PayloadConnect) Reset()                    { *m = PayloadConnect{} }
func (m *PayloadConnect) String() string            { return proto.CompactTextString(m) }
func (*PayloadConnect) ProtoMessage()               {}
func (*PayloadConnect) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *PayloadConnect) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

type PayloadDisconnect struct {
	// id is the identifier of the node that should be disconnected
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *PayloadDisconnect) Reset()                    { *m = PayloadDisconnect{} }
func (m *PayloadDisconnect) String() string            { return proto.CompactTextString(m) }
func (*PayloadDisconnect) ProtoMessage()               {}
func (*PayloadDisconnect) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *PayloadDisconnect) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type PayloadNode struct {
	// id is the identifier of the node. Usually a sha1 hash.
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// baddr is the booster listening address.
	Baddr string `protobuf:"bytes,2,opt,name=baddr" json:"baddr,omitempty"`
	// paddr is the proxy listening address.
	Paddr string `protobuf:"bytes,3,opt,name=paddr" json:"paddr,omitempty"`
	// active tells the connection state of the node.
	Active bool `protobuf:"varint,4,opt,name=active" json:"active,omitempty"`
	// tunnels are the proxy tunnels managed by this node.
	Tunnels []*PayloadNode_Tunnel `protobuf:"bytes,5,rep,name=tunnels" json:"tunnels,omitempty"`
}

func (m *PayloadNode) Reset()                    { *m = PayloadNode{} }
func (m *PayloadNode) String() string            { return proto.CompactTextString(m) }
func (*PayloadNode) ProtoMessage()               {}
func (*PayloadNode) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *PayloadNode) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PayloadNode) GetBaddr() string {
	if m != nil {
		return m.Baddr
	}
	return ""
}

func (m *PayloadNode) GetPaddr() string {
	if m != nil {
		return m.Paddr
	}
	return ""
}

func (m *PayloadNode) GetActive() bool {
	if m != nil {
		return m.Active
	}
	return false
}

func (m *PayloadNode) GetTunnels() []*PayloadNode_Tunnel {
	if m != nil {
		return m.Tunnels
	}
	return nil
}

type PayloadNode_Tunnel struct {
	// id is the tunnel identifier. Usally a sha1 hash.
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// target is the remote endpoint address of the tunnel.
	Target string `protobuf:"bytes,2,opt,name=target" json:"target,omitempty"`
	// acks is the number of acknoledgments on this tunnel.
	Acks int32 `protobuf:"varint,3,opt,name=acks" json:"acks,omitempty"`
	// copies are the replications of this tunnel.
	Copies int32 `protobuf:"varint,4,opt,name=copies" json:"copies,omitempty"`
}

func (m *PayloadNode_Tunnel) Reset()                    { *m = PayloadNode_Tunnel{} }
func (m *PayloadNode_Tunnel) String() string            { return proto.CompactTextString(m) }
func (*PayloadNode_Tunnel) ProtoMessage()               {}
func (*PayloadNode_Tunnel) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6, 0} }

func (m *PayloadNode_Tunnel) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PayloadNode_Tunnel) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *PayloadNode_Tunnel) GetAcks() int32 {
	if m != nil {
		return m.Acks
	}
	return 0
}

func (m *PayloadNode_Tunnel) GetCopies() int32 {
	if m != nil {
		return m.Copies
	}
	return 0
}

type PayloadHeartbeat struct {
	// id is the identifier of the heartbeat message. Should be unique.
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// hops is the number of times that the heartbeat message has been reused.
	Hops int32 `protobuf:"varint,2,opt,name=hops" json:"hops,omitempty"`
	// ttl is the time to leave.
	Ttl *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=ttl" json:"ttl,omitempty"`
}

func (m *PayloadHeartbeat) Reset()                    { *m = PayloadHeartbeat{} }
func (m *PayloadHeartbeat) String() string            { return proto.CompactTextString(m) }
func (*PayloadHeartbeat) ProtoMessage()               {}
func (*PayloadHeartbeat) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *PayloadHeartbeat) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PayloadHeartbeat) GetHops() int32 {
	if m != nil {
		return m.Hops
	}
	return 0
}

func (m *PayloadHeartbeat) GetTtl() *google_protobuf.Timestamp {
	if m != nil {
		return m.Ttl
	}
	return nil
}

type PayloadTunnelEvent struct {
	// target is the remote endpoint address of the tunnel.
	Target string `protobuf:"bytes,1,opt,name=target" json:"target,omitempty"`
	// event is the operation performed on the tunnel.
	Event int32 `protobuf:"varint,2,opt,name=event" json:"event,omitempty"`
}

func (m *PayloadTunnelEvent) Reset()                    { *m = PayloadTunnelEvent{} }
func (m *PayloadTunnelEvent) String() string            { return proto.CompactTextString(m) }
func (*PayloadTunnelEvent) ProtoMessage()               {}
func (*PayloadTunnelEvent) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

func (m *PayloadTunnelEvent) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *PayloadTunnelEvent) GetEvent() int32 {
	if m != nil {
		return m.Event
	}
	return 0
}

func init() {
	proto.RegisterType((*PayloadCtrl)(nil), "internal.PayloadCtrl")
	proto.RegisterType((*PayloadBandwidth)(nil), "internal.PayloadBandwidth")
	proto.RegisterType((*PayloadInspect)(nil), "internal.PayloadInspect")
	proto.RegisterType((*PayloadHello)(nil), "internal.PayloadHello")
	proto.RegisterType((*PayloadConnect)(nil), "internal.PayloadConnect")
	proto.RegisterType((*PayloadDisconnect)(nil), "internal.PayloadDisconnect")
	proto.RegisterType((*PayloadNode)(nil), "internal.PayloadNode")
	proto.RegisterType((*PayloadNode_Tunnel)(nil), "internal.PayloadNode.Tunnel")
	proto.RegisterType((*PayloadHeartbeat)(nil), "internal.PayloadHeartbeat")
	proto.RegisterType((*PayloadTunnelEvent)(nil), "internal.PayloadTunnelEvent")
}

func init() { proto.RegisterFile("payload.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 433 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0xc1, 0x6a, 0xdc, 0x30,
	0x10, 0xc5, 0x76, 0xbc, 0xcd, 0xce, 0xb6, 0x21, 0x15, 0x25, 0x98, 0x25, 0x50, 0xe3, 0x5e, 0x0c,
	0x0d, 0x0e, 0xa4, 0xd0, 0x43, 0x8f, 0x49, 0x0b, 0xe9, 0xa5, 0x14, 0x11, 0x7a, 0xea, 0x45, 0xb6,
	0x27, 0x1b, 0x53, 0x47, 0x12, 0xf2, 0x6c, 0x4a, 0xbe, 0xbc, 0xd7, 0xa2, 0x91, 0xbc, 0x1b, 0x1a,
	0x7a, 0x9b, 0xf7, 0xf4, 0x34, 0xef, 0x49, 0x33, 0xf0, 0xca, 0xaa, 0xc7, 0xd1, 0xa8, 0xbe, 0xb1,
	0xce, 0x90, 0x11, 0x87, 0x83, 0x26, 0x74, 0x5a, 0x8d, 0xeb, 0xb7, 0x1b, 0x63, 0x36, 0x23, 0x9e,
	0x33, 0xdf, 0x6e, 0x6f, 0xcf, 0x69, 0xb8, 0xc7, 0x89, 0xd4, 0xbd, 0x0d, 0xd2, 0xea, 0x3d, 0xac,
	0xbe, 0x87, 0xbb, 0x57, 0xe4, 0x46, 0x71, 0x0a, 0x4b, 0x63, 0xd1, 0x29, 0x1a, 0x8c, 0x2e, 0x92,
	0x32, 0xa9, 0x73, 0xb9, 0x27, 0xaa, 0x1f, 0x70, 0x1c, 0xc5, 0x97, 0x4a, 0xf7, 0xbf, 0x87, 0x9e,
	0xee, 0xc4, 0x31, 0x64, 0x64, 0x88, 0xb5, 0x99, 0xf4, 0xa5, 0xef, 0xd1, 0xce, 0xc7, 0x45, 0xca,
	0xfc, 0x9e, 0x10, 0x02, 0x0e, 0xe8, 0xd1, 0x62, 0x91, 0x95, 0x49, 0xbd, 0x94, 0x5c, 0x57, 0x67,
	0x70, 0x14, 0xfb, 0x7e, 0xd5, 0x93, 0xc5, 0x8e, 0xc4, 0x1a, 0x0e, 0x6f, 0x51, 0xd1, 0xd6, 0xe1,
	0x54, 0x24, 0x65, 0x56, 0xe7, 0x72, 0x87, 0xab, 0x4f, 0xf0, 0x32, 0xaa, 0xaf, 0x71, 0x1c, 0x8d,
	0x78, 0x03, 0x79, 0x6b, 0x8d, 0x0b, 0x19, 0x96, 0x32, 0x00, 0xcf, 0x5a, 0x66, 0xd3, 0xc0, 0x32,
	0xa8, 0xea, 0x9d, 0xd3, 0x95, 0xd1, 0xda, 0x3b, 0x9d, 0xc0, 0x82, 0x94, 0xdb, 0xe0, 0x7c, 0x3d,
	0xa2, 0xea, 0x1d, 0xbc, 0x8e, 0xca, 0xcf, 0xc3, 0xd4, 0x45, 0xf1, 0x11, 0xa4, 0x43, 0x1f, 0x85,
	0xe9, 0xd0, 0x57, 0x7f, 0x92, 0xdd, 0xf7, 0x7d, 0x33, 0x3d, 0xfe, 0x7b, 0xce, 0xd1, 0x54, 0xdf,
	0xbb, 0x39, 0x04, 0x03, 0x8e, 0xc6, 0x6c, 0x16, 0xa3, 0x31, 0x7b, 0x02, 0x0b, 0xd5, 0xd1, 0xf0,
	0x80, 0xc5, 0x41, 0x99, 0xd4, 0x87, 0x32, 0x22, 0xf1, 0x11, 0x5e, 0xd0, 0x56, 0x6b, 0x1c, 0xa7,
	0x22, 0x2f, 0xb3, 0x7a, 0x75, 0x71, 0xda, 0xcc, 0xe3, 0x6d, 0x9e, 0x78, 0x37, 0x37, 0x2c, 0x92,
	0xb3, 0x78, 0xfd, 0x13, 0x16, 0x81, 0x7a, 0x96, 0x6a, 0xff, 0xe4, 0xf4, 0xe9, 0x93, 0xfd, 0x68,
	0x54, 0xf7, 0x6b, 0xe2, 0x58, 0xb9, 0xe4, 0xda, 0x6b, 0x3b, 0x63, 0x07, 0x9c, 0x38, 0x55, 0x2e,
	0x23, 0xaa, 0xfa, 0xdd, 0x2a, 0x5c, 0xa3, 0x72, 0xd4, 0xa2, 0x7a, 0xf6, 0x3b, 0xbe, 0xdf, 0x9d,
	0xb1, 0x13, 0xbb, 0xe4, 0x92, 0x6b, 0x71, 0x06, 0x19, 0xd1, 0xc8, 0x16, 0xab, 0x8b, 0x75, 0x13,
	0xd6, 0xb3, 0x99, 0xd7, 0xb3, 0xb9, 0x99, 0xd7, 0x53, 0x7a, 0x59, 0x75, 0x09, 0x22, 0xba, 0x84,
	0xa7, 0x7c, 0x79, 0x40, 0xfd, 0xdf, 0x91, 0xf9, 0x7f, 0x45, 0x2f, 0x88, 0x86, 0x01, 0xb4, 0x0b,
	0x6e, 0xfe, 0xe1, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x55, 0x8b, 0x0c, 0x24, 0x03, 0x00,
	0x00,
}
