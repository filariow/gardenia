// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: rosina.proto

package valvedprotos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RosinaSvcClient is the client API for RosinaSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RosinaSvcClient interface {
	OpenValve(ctx context.Context, in *OpenRequest, opts ...grpc.CallOption) (*OpenReply, error)
	CloseValve(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseReply, error)
}

type rosinaSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewRosinaSvcClient(cc grpc.ClientConnInterface) RosinaSvcClient {
	return &rosinaSvcClient{cc}
}

func (c *rosinaSvcClient) OpenValve(ctx context.Context, in *OpenRequest, opts ...grpc.CallOption) (*OpenReply, error) {
	out := new(OpenReply)
	err := c.cc.Invoke(ctx, "/valvedgrpc.RosinaSvc/OpenValve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rosinaSvcClient) CloseValve(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseReply, error) {
	out := new(CloseReply)
	err := c.cc.Invoke(ctx, "/valvedgrpc.RosinaSvc/CloseValve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RosinaSvcServer is the server API for RosinaSvc service.
// All implementations must embed UnimplementedRosinaSvcServer
// for forward compatibility
type RosinaSvcServer interface {
	OpenValve(context.Context, *OpenRequest) (*OpenReply, error)
	CloseValve(context.Context, *CloseRequest) (*CloseReply, error)
	mustEmbedUnimplementedRosinaSvcServer()
}

// UnimplementedRosinaSvcServer must be embedded to have forward compatible implementations.
type UnimplementedRosinaSvcServer struct {
}

func (UnimplementedRosinaSvcServer) OpenValve(context.Context, *OpenRequest) (*OpenReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenValve not implemented")
}
func (UnimplementedRosinaSvcServer) CloseValve(context.Context, *CloseRequest) (*CloseReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseValve not implemented")
}
func (UnimplementedRosinaSvcServer) mustEmbedUnimplementedRosinaSvcServer() {}

// UnsafeRosinaSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RosinaSvcServer will
// result in compilation errors.
type UnsafeRosinaSvcServer interface {
	mustEmbedUnimplementedRosinaSvcServer()
}

func RegisterRosinaSvcServer(s grpc.ServiceRegistrar, srv RosinaSvcServer) {
	s.RegisterService(&RosinaSvc_ServiceDesc, srv)
}

func _RosinaSvc_OpenValve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RosinaSvcServer).OpenValve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/valvedgrpc.RosinaSvc/OpenValve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RosinaSvcServer).OpenValve(ctx, req.(*OpenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RosinaSvc_CloseValve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RosinaSvcServer).CloseValve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/valvedgrpc.RosinaSvc/CloseValve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RosinaSvcServer).CloseValve(ctx, req.(*CloseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RosinaSvc_ServiceDesc is the grpc.ServiceDesc for RosinaSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RosinaSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "valvedgrpc.RosinaSvc",
	HandlerType: (*RosinaSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OpenValve",
			Handler:    _RosinaSvc_OpenValve_Handler,
		},
		{
			MethodName: "CloseValve",
			Handler:    _RosinaSvc_CloseValve_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rosina.proto",
}
