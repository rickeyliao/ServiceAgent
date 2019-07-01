// Code generated by protoc-gen-go. DO NOT EDIT.
// source: licenseuser.proto

package rpccmdservice

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type LicenseUserChgReq struct {
	Op                   bool     `protobuf:"varint,1,opt,name=op,proto3" json:"op,omitempty"`
	User                 string   `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LicenseUserChgReq) Reset()         { *m = LicenseUserChgReq{} }
func (m *LicenseUserChgReq) String() string { return proto.CompactTextString(m) }
func (*LicenseUserChgReq) ProtoMessage()    {}
func (*LicenseUserChgReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_5faa6abdc8a457b0, []int{0}
}

func (m *LicenseUserChgReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LicenseUserChgReq.Unmarshal(m, b)
}
func (m *LicenseUserChgReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LicenseUserChgReq.Marshal(b, m, deterministic)
}
func (m *LicenseUserChgReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LicenseUserChgReq.Merge(m, src)
}
func (m *LicenseUserChgReq) XXX_Size() int {
	return xxx_messageInfo_LicenseUserChgReq.Size(m)
}
func (m *LicenseUserChgReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LicenseUserChgReq.DiscardUnknown(m)
}

var xxx_messageInfo_LicenseUserChgReq proto.InternalMessageInfo

func (m *LicenseUserChgReq) GetOp() bool {
	if m != nil {
		return m.Op
	}
	return false
}

func (m *LicenseUserChgReq) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func init() {
	proto.RegisterType((*LicenseUserChgReq)(nil), "rpccmdservice.LicenseUserChgReq")
}

func init() { proto.RegisterFile("licenseuser.proto", fileDescriptor_5faa6abdc8a457b0) }

var fileDescriptor_5faa6abdc8a457b0 = []byte{
	// 159 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xcc, 0xc9, 0x4c, 0x4e,
	0xcd, 0x2b, 0x4e, 0x2d, 0x2d, 0x4e, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2d,
	0x2a, 0x48, 0x4e, 0xce, 0x4d, 0x29, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x95, 0xe2, 0x4d, 0x49,
	0x4d, 0x4b, 0x2c, 0xcd, 0x29, 0x81, 0xc8, 0x2a, 0x99, 0x73, 0x09, 0xfa, 0x40, 0xb4, 0x84, 0x16,
	0xa7, 0x16, 0x39, 0x67, 0xa4, 0x07, 0xa5, 0x16, 0x0a, 0xf1, 0x71, 0x31, 0xe5, 0x17, 0x48, 0x30,
	0x2a, 0x30, 0x6a, 0x70, 0x04, 0x31, 0xe5, 0x17, 0x08, 0x09, 0x71, 0xb1, 0x80, 0x0c, 0x94, 0x60,
	0x52, 0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb3, 0x8d, 0x92, 0xb8, 0xf8, 0x50, 0x35, 0x0a, 0x05, 0x70,
	0xf1, 0x39, 0x67, 0xa4, 0x23, 0x09, 0x0a, 0x29, 0xe8, 0xa1, 0xd8, 0xad, 0x87, 0x61, 0x93, 0x94,
	0x14, 0x9a, 0x0a, 0x17, 0x88, 0xe3, 0x82, 0x52, 0x8b, 0x0b, 0x94, 0x18, 0x92, 0xd8, 0xc0, 0x6e,
	0x34, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xa7, 0xcf, 0xdf, 0x1a, 0xd6, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LicenseUserChgClient is the client API for LicenseUserChg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LicenseUserChgClient interface {
	ChgLicenseUser(ctx context.Context, in *LicenseUserChgReq, opts ...grpc.CallOption) (*DefaultResp, error)
}

type licenseUserChgClient struct {
	cc *grpc.ClientConn
}

func NewLicenseUserChgClient(cc *grpc.ClientConn) LicenseUserChgClient {
	return &licenseUserChgClient{cc}
}

func (c *licenseUserChgClient) ChgLicenseUser(ctx context.Context, in *LicenseUserChgReq, opts ...grpc.CallOption) (*DefaultResp, error) {
	out := new(DefaultResp)
	err := c.cc.Invoke(ctx, "/rpccmdservice.LicenseUserChg/ChgLicenseUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LicenseUserChgServer is the server API for LicenseUserChg service.
type LicenseUserChgServer interface {
	ChgLicenseUser(context.Context, *LicenseUserChgReq) (*DefaultResp, error)
}

func RegisterLicenseUserChgServer(s *grpc.Server, srv LicenseUserChgServer) {
	s.RegisterService(&_LicenseUserChg_serviceDesc, srv)
}

func _LicenseUserChg_ChgLicenseUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LicenseUserChgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LicenseUserChgServer).ChgLicenseUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpccmdservice.LicenseUserChg/ChgLicenseUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LicenseUserChgServer).ChgLicenseUser(ctx, req.(*LicenseUserChgReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _LicenseUserChg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpccmdservice.LicenseUserChg",
	HandlerType: (*LicenseUserChgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ChgLicenseUser",
			Handler:    _LicenseUserChg_ChgLicenseUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "licenseuser.proto",
}
