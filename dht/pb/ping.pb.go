// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ping.proto

package pbdht

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Pingreq struct {
	Sn                   uint64   `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	Msgtyp               uint32   `protobuf:"varint,2,opt,name=msgtyp,proto3" json:"msgtyp,omitempty"`
	Nbsaddr              []byte   `protobuf:"bytes,3,opt,name=nbsaddr,proto3" json:"nbsaddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pingreq) Reset()         { *m = Pingreq{} }
func (m *Pingreq) String() string { return proto.CompactTextString(m) }
func (*Pingreq) ProtoMessage()    {}
func (*Pingreq) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d51d96c3ad891f5, []int{0}
}

func (m *Pingreq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pingreq.Unmarshal(m, b)
}
func (m *Pingreq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pingreq.Marshal(b, m, deterministic)
}
func (m *Pingreq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pingreq.Merge(m, src)
}
func (m *Pingreq) XXX_Size() int {
	return xxx_messageInfo_Pingreq.Size(m)
}
func (m *Pingreq) XXX_DiscardUnknown() {
	xxx_messageInfo_Pingreq.DiscardUnknown(m)
}

var xxx_messageInfo_Pingreq proto.InternalMessageInfo

func (m *Pingreq) GetSn() uint64 {
	if m != nil {
		return m.Sn
	}
	return 0
}

func (m *Pingreq) GetMsgtyp() uint32 {
	if m != nil {
		return m.Msgtyp
	}
	return 0
}

func (m *Pingreq) GetNbsaddr() []byte {
	if m != nil {
		return m.Nbsaddr
	}
	return nil
}

type Pingresp struct {
	Rcvsn                uint64   `protobuf:"varint,1,opt,name=rcvsn,proto3" json:"rcvsn,omitempty"`
	Msgtype              uint32   `protobuf:"varint,2,opt,name=msgtype,proto3" json:"msgtype,omitempty"`
	Nbsaddr              []byte   `protobuf:"bytes,3,opt,name=nbsaddr,proto3" json:"nbsaddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pingresp) Reset()         { *m = Pingresp{} }
func (m *Pingresp) String() string { return proto.CompactTextString(m) }
func (*Pingresp) ProtoMessage()    {}
func (*Pingresp) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d51d96c3ad891f5, []int{1}
}

func (m *Pingresp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pingresp.Unmarshal(m, b)
}
func (m *Pingresp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pingresp.Marshal(b, m, deterministic)
}
func (m *Pingresp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pingresp.Merge(m, src)
}
func (m *Pingresp) XXX_Size() int {
	return xxx_messageInfo_Pingresp.Size(m)
}
func (m *Pingresp) XXX_DiscardUnknown() {
	xxx_messageInfo_Pingresp.DiscardUnknown(m)
}

var xxx_messageInfo_Pingresp proto.InternalMessageInfo

func (m *Pingresp) GetRcvsn() uint64 {
	if m != nil {
		return m.Rcvsn
	}
	return 0
}

func (m *Pingresp) GetMsgtype() uint32 {
	if m != nil {
		return m.Msgtype
	}
	return 0
}

func (m *Pingresp) GetNbsaddr() []byte {
	if m != nil {
		return m.Nbsaddr
	}
	return nil
}

func init() {
	proto.RegisterType((*Pingreq)(nil), "pbdht.pingreq")
	proto.RegisterType((*Pingresp)(nil), "pbdht.pingresp")
}

func init() { proto.RegisterFile("ping.proto", fileDescriptor_6d51d96c3ad891f5) }

var fileDescriptor_6d51d96c3ad891f5 = []byte{
	// 142 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0xc8, 0xcc, 0x4b,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0x48, 0x4a, 0xc9, 0x28, 0x51, 0xf2, 0xe6,
	0x62, 0x07, 0x09, 0x16, 0xa5, 0x16, 0x0a, 0xf1, 0x71, 0x31, 0x15, 0xe7, 0x49, 0x30, 0x2a, 0x30,
	0x6a, 0xb0, 0x04, 0x31, 0x15, 0xe7, 0x09, 0x89, 0x71, 0xb1, 0xe5, 0x16, 0xa7, 0x97, 0x54, 0x16,
	0x48, 0x30, 0x29, 0x30, 0x6a, 0xf0, 0x06, 0x41, 0x79, 0x42, 0x12, 0x5c, 0xec, 0x79, 0x49, 0xc5,
	0x89, 0x29, 0x29, 0x45, 0x12, 0xcc, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x30, 0xae, 0x52, 0x08, 0x17,
	0x07, 0xc4, 0xb0, 0xe2, 0x02, 0x21, 0x11, 0x2e, 0xd6, 0xa2, 0xe4, 0x32, 0xb8, 0x81, 0x10, 0x0e,
	0x48, 0x2f, 0xc4, 0x94, 0x54, 0xa8, 0xa1, 0x30, 0x2e, 0x6e, 0x53, 0x93, 0xd8, 0xc0, 0x0e, 0x36,
	0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xa2, 0xa1, 0xef, 0x87, 0xbe, 0x00, 0x00, 0x00,
}