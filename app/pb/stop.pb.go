// Code generated by protoc-gen-go. DO NOT EDIT.
// source: stop.proto

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

func init() { proto.RegisterFile("stop.proto", fileDescriptor_f049a61f03aafc0b) }

var fileDescriptor_f049a61f03aafc0b = []byte{
	// 114 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x2e, 0xc9, 0x2f,
	0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2d, 0x2a, 0x48, 0x4e, 0xce, 0x4d, 0x29, 0x4e,
	0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x95, 0xe2, 0x4d, 0x49, 0x4d, 0x4b, 0x2c, 0xcd, 0x29, 0x81, 0xc8,
	0x1a, 0xc5, 0x71, 0xf1, 0x43, 0x05, 0xf2, 0x92, 0x8a, 0x8b, 0x13, 0x8b, 0x8b, 0xca, 0x84, 0xbc,
	0xb9, 0x78, 0x5c, 0x20, 0x42, 0x7e, 0x20, 0x21, 0x21, 0x59, 0x3d, 0x14, 0x13, 0xf4, 0xa0, 0x92,
	0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x52, 0x52, 0xb8, 0xa4, 0x8b, 0x0b, 0x94, 0x18, 0x92,
	0xd8, 0xc0, 0xd6, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x3e, 0x8e, 0x59, 0xa4, 0x92, 0x00,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DefaultnbssasrvClient is the client API for Defaultnbssasrv service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DefaultnbssasrvClient interface {
	DefaultNbssa(ctx context.Context, in *DefaultRequest, opts ...grpc.CallOption) (*DefaultResp, error)
}

type defaultnbssasrvClient struct {
	cc *grpc.ClientConn
}

func NewDefaultnbssasrvClient(cc *grpc.ClientConn) DefaultnbssasrvClient {
	return &defaultnbssasrvClient{cc}
}

func (c *defaultnbssasrvClient) DefaultNbssa(ctx context.Context, in *DefaultRequest, opts ...grpc.CallOption) (*DefaultResp, error) {
	out := new(DefaultResp)
	err := c.cc.Invoke(ctx, "/rpccmdservice.defaultnbssasrv/DefaultNbssa", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DefaultnbssasrvServer is the server API for Defaultnbssasrv service.
type DefaultnbssasrvServer interface {
	DefaultNbssa(context.Context, *DefaultRequest) (*DefaultResp, error)
}

func RegisterDefaultnbssasrvServer(s *grpc.Server, srv DefaultnbssasrvServer) {
	s.RegisterService(&_Defaultnbssasrv_serviceDesc, srv)
}

func _Defaultnbssasrv_DefaultNbssa_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DefaultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DefaultnbssasrvServer).DefaultNbssa(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpccmdservice.defaultnbssasrv/DefaultNbssa",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DefaultnbssasrvServer).DefaultNbssa(ctx, req.(*DefaultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Defaultnbssasrv_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpccmdservice.defaultnbssasrv",
	HandlerType: (*DefaultnbssasrvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DefaultNbssa",
			Handler:    _Defaultnbssasrv_DefaultNbssa_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stop.proto",
}
