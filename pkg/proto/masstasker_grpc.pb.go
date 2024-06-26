// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: masstasker.proto

package masstasker

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

const (
	MassTasker_Update_FullMethodName  = "/masstasker.MassTasker/Update"
	MassTasker_Query_FullMethodName   = "/masstasker.MassTasker/Query"
	MassTasker_BulkSet_FullMethodName = "/masstasker.MassTasker/BulkSet"
	MassTasker_Debug_FullMethodName   = "/masstasker.MassTasker/Debug"
)

// MassTaskerClient is the client API for MassTasker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MassTaskerClient interface {
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error)
	BulkSet(ctx context.Context, in *BulkSetRequest, opts ...grpc.CallOption) (*BulkSetResponse, error)
	Debug(ctx context.Context, in *DebugRequest, opts ...grpc.CallOption) (*DebugResponse, error)
}

type massTaskerClient struct {
	cc grpc.ClientConnInterface
}

func NewMassTaskerClient(cc grpc.ClientConnInterface) MassTaskerClient {
	return &massTaskerClient{cc}
}

func (c *massTaskerClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, MassTasker_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *massTaskerClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error) {
	out := new(QueryResponse)
	err := c.cc.Invoke(ctx, MassTasker_Query_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *massTaskerClient) BulkSet(ctx context.Context, in *BulkSetRequest, opts ...grpc.CallOption) (*BulkSetResponse, error) {
	out := new(BulkSetResponse)
	err := c.cc.Invoke(ctx, MassTasker_BulkSet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *massTaskerClient) Debug(ctx context.Context, in *DebugRequest, opts ...grpc.CallOption) (*DebugResponse, error) {
	out := new(DebugResponse)
	err := c.cc.Invoke(ctx, MassTasker_Debug_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MassTaskerServer is the server API for MassTasker service.
// All implementations must embed UnimplementedMassTaskerServer
// for forward compatibility
type MassTaskerServer interface {
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	Query(context.Context, *QueryRequest) (*QueryResponse, error)
	BulkSet(context.Context, *BulkSetRequest) (*BulkSetResponse, error)
	Debug(context.Context, *DebugRequest) (*DebugResponse, error)
	mustEmbedUnimplementedMassTaskerServer()
}

// UnimplementedMassTaskerServer must be embedded to have forward compatible implementations.
type UnimplementedMassTaskerServer struct {
}

func (UnimplementedMassTaskerServer) Update(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedMassTaskerServer) Query(context.Context, *QueryRequest) (*QueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedMassTaskerServer) BulkSet(context.Context, *BulkSetRequest) (*BulkSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BulkSet not implemented")
}
func (UnimplementedMassTaskerServer) Debug(context.Context, *DebugRequest) (*DebugResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Debug not implemented")
}
func (UnimplementedMassTaskerServer) mustEmbedUnimplementedMassTaskerServer() {}

// UnsafeMassTaskerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MassTaskerServer will
// result in compilation errors.
type UnsafeMassTaskerServer interface {
	mustEmbedUnimplementedMassTaskerServer()
}

func RegisterMassTaskerServer(s grpc.ServiceRegistrar, srv MassTaskerServer) {
	s.RegisterService(&MassTasker_ServiceDesc, srv)
}

func _MassTasker_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MassTaskerServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MassTasker_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MassTaskerServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MassTasker_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MassTaskerServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MassTasker_Query_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MassTaskerServer).Query(ctx, req.(*QueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MassTasker_BulkSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MassTaskerServer).BulkSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MassTasker_BulkSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MassTaskerServer).BulkSet(ctx, req.(*BulkSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MassTasker_Debug_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DebugRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MassTaskerServer).Debug(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MassTasker_Debug_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MassTaskerServer).Debug(ctx, req.(*DebugRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MassTasker_ServiceDesc is the grpc.ServiceDesc for MassTasker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MassTasker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "masstasker.MassTasker",
	HandlerType: (*MassTaskerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _MassTasker_Update_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _MassTasker_Query_Handler,
		},
		{
			MethodName: "BulkSet",
			Handler:    _MassTasker_BulkSet_Handler,
		},
		{
			MethodName: "Debug",
			Handler:    _MassTasker_Debug_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "masstasker.proto",
}
