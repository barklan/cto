// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package porter

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

// PorterClient is the client API for Porter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PorterClient interface {
	ProjectAlert(ctx context.Context, in *ProjectAlertRequest, opts ...grpc.CallOption) (*Message, error)
	InternalAlert(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
	NewIssue(ctx context.Context, in *NewIssueRequest, opts ...grpc.CallOption) (*Message, error)
}

type porterClient struct {
	cc grpc.ClientConnInterface
}

func NewPorterClient(cc grpc.ClientConnInterface) PorterClient {
	return &porterClient{cc}
}

func (c *porterClient) ProjectAlert(ctx context.Context, in *ProjectAlertRequest, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/protos.Porter/ProjectAlert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *porterClient) InternalAlert(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/protos.Porter/InternalAlert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *porterClient) NewIssue(ctx context.Context, in *NewIssueRequest, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/protos.Porter/NewIssue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PorterServer is the server API for Porter service.
// All implementations must embed UnimplementedPorterServer
// for forward compatibility
type PorterServer interface {
	ProjectAlert(context.Context, *ProjectAlertRequest) (*Message, error)
	InternalAlert(context.Context, *Message) (*Message, error)
	NewIssue(context.Context, *NewIssueRequest) (*Message, error)
	mustEmbedUnimplementedPorterServer()
}

// UnimplementedPorterServer must be embedded to have forward compatible implementations.
type UnimplementedPorterServer struct {
}

func (UnimplementedPorterServer) ProjectAlert(context.Context, *ProjectAlertRequest) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProjectAlert not implemented")
}
func (UnimplementedPorterServer) InternalAlert(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InternalAlert not implemented")
}
func (UnimplementedPorterServer) NewIssue(context.Context, *NewIssueRequest) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewIssue not implemented")
}
func (UnimplementedPorterServer) mustEmbedUnimplementedPorterServer() {}

// UnsafePorterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PorterServer will
// result in compilation errors.
type UnsafePorterServer interface {
	mustEmbedUnimplementedPorterServer()
}

func RegisterPorterServer(s grpc.ServiceRegistrar, srv PorterServer) {
	s.RegisterService(&Porter_ServiceDesc, srv)
}

func _Porter_ProjectAlert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectAlertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PorterServer).ProjectAlert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Porter/ProjectAlert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PorterServer).ProjectAlert(ctx, req.(*ProjectAlertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Porter_InternalAlert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PorterServer).InternalAlert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Porter/InternalAlert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PorterServer).InternalAlert(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _Porter_NewIssue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewIssueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PorterServer).NewIssue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Porter/NewIssue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PorterServer).NewIssue(ctx, req.(*NewIssueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Porter_ServiceDesc is the grpc.ServiceDesc for Porter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Porter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Porter",
	HandlerType: (*PorterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProjectAlert",
			Handler:    _Porter_ProjectAlert_Handler,
		},
		{
			MethodName: "InternalAlert",
			Handler:    _Porter_InternalAlert_Handler,
		},
		{
			MethodName: "NewIssue",
			Handler:    _Porter_NewIssue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/protos/porter/main.proto",
}
