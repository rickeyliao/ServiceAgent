// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dhtstore.proto

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

type Dhtstore struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Share                bool     `protobuf:"varint,2,opt,name=share,proto3" json:"share,omitempty"`
	Value                [][]byte `protobuf:"bytes,3,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Dhtstore) Reset()         { *m = Dhtstore{} }
func (m *Dhtstore) String() string { return proto.CompactTextString(m) }
func (*Dhtstore) ProtoMessage()    {}
func (*Dhtstore) Descriptor() ([]byte, []int) {
	return fileDescriptor_e022415cdf0545ce, []int{0}
}

func (m *Dhtstore) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Dhtstore.Unmarshal(m, b)
}
func (m *Dhtstore) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Dhtstore.Marshal(b, m, deterministic)
}
func (m *Dhtstore) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Dhtstore.Merge(m, src)
}
func (m *Dhtstore) XXX_Size() int {
	return xxx_messageInfo_Dhtstore.Size(m)
}
func (m *Dhtstore) XXX_DiscardUnknown() {
	xxx_messageInfo_Dhtstore.DiscardUnknown(m)
}

var xxx_messageInfo_Dhtstore proto.InternalMessageInfo

func (m *Dhtstore) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Dhtstore) GetShare() bool {
	if m != nil {
		return m.Share
	}
	return false
}

func (m *Dhtstore) GetValue() [][]byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*Dhtstore)(nil), "pbdht.dhtstore")
}

func init() {
	proto.RegisterFile("dhtstore.proto", fileDescriptor_e022415cdf0545ce)
}

var fileDescriptor_e022415cdf0545ce = []byte{
	// 107 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0xc9, 0x28, 0x29,
	0x2e, 0xc9, 0x2f, 0x4a, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0x48, 0x4a, 0xc9,
	0x28, 0x51, 0xf2, 0xe0, 0xe2, 0x80, 0x49, 0x08, 0x09, 0x70, 0x31, 0x67, 0xa7, 0x56, 0x4a, 0x30,
	0x2a, 0x30, 0x6a, 0xf0, 0x04, 0x81, 0x98, 0x42, 0x22, 0x5c, 0xac, 0xc5, 0x19, 0x89, 0x45, 0xa9,
	0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x1c, 0x41, 0x10, 0x0e, 0x48, 0xb4, 0x2c, 0x31, 0xa7, 0x34, 0x55,
	0x82, 0x59, 0x81, 0x59, 0x83, 0x27, 0x08, 0xc2, 0x49, 0x62, 0x03, 0x9b, 0x6b, 0x0c, 0x08, 0x00,
	0x00, 0xff, 0xff, 0xe5, 0x92, 0x4e, 0xe9, 0x69, 0x00, 0x00, 0x00,
}
