// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.26.1
// source: helloworld/v1/cell.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Cell_CreateCell_FullMethodName = "/api.helloworld.v1.Cell/CreateCell"
	Cell_UpdateCell_FullMethodName = "/api.helloworld.v1.Cell/UpdateCell"
	Cell_DeleteCell_FullMethodName = "/api.helloworld.v1.Cell/DeleteCell"
	Cell_GetCell_FullMethodName    = "/api.helloworld.v1.Cell/GetCell"
	Cell_ListCell_FullMethodName   = "/api.helloworld.v1.Cell/ListCell"
)

// CellClient is the client API for Cell service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CellClient interface {
	CreateCell(ctx context.Context, in *CreateCellRequest, opts ...grpc.CallOption) (*CreateCellReply, error)
	UpdateCell(ctx context.Context, in *UpdateCellRequest, opts ...grpc.CallOption) (*UpdateCellReply, error)
	DeleteCell(ctx context.Context, in *DeleteCellRequest, opts ...grpc.CallOption) (*DeleteCellReply, error)
	GetCell(ctx context.Context, in *GetCellRequest, opts ...grpc.CallOption) (*GetCellReply, error)
	ListCell(ctx context.Context, in *ListCellRequest, opts ...grpc.CallOption) (*ListCellReply, error)
}

type cellClient struct {
	cc grpc.ClientConnInterface
}

func NewCellClient(cc grpc.ClientConnInterface) CellClient {
	return &cellClient{cc}
}

func (c *cellClient) CreateCell(ctx context.Context, in *CreateCellRequest, opts ...grpc.CallOption) (*CreateCellReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCellReply)
	err := c.cc.Invoke(ctx, Cell_CreateCell_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cellClient) UpdateCell(ctx context.Context, in *UpdateCellRequest, opts ...grpc.CallOption) (*UpdateCellReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateCellReply)
	err := c.cc.Invoke(ctx, Cell_UpdateCell_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cellClient) DeleteCell(ctx context.Context, in *DeleteCellRequest, opts ...grpc.CallOption) (*DeleteCellReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteCellReply)
	err := c.cc.Invoke(ctx, Cell_DeleteCell_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cellClient) GetCell(ctx context.Context, in *GetCellRequest, opts ...grpc.CallOption) (*GetCellReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCellReply)
	err := c.cc.Invoke(ctx, Cell_GetCell_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cellClient) ListCell(ctx context.Context, in *ListCellRequest, opts ...grpc.CallOption) (*ListCellReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCellReply)
	err := c.cc.Invoke(ctx, Cell_ListCell_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CellServer is the server API for Cell service.
// All implementations must embed UnimplementedCellServer
// for forward compatibility.
type CellServer interface {
	CreateCell(context.Context, *CreateCellRequest) (*CreateCellReply, error)
	UpdateCell(context.Context, *UpdateCellRequest) (*UpdateCellReply, error)
	DeleteCell(context.Context, *DeleteCellRequest) (*DeleteCellReply, error)
	GetCell(context.Context, *GetCellRequest) (*GetCellReply, error)
	ListCell(context.Context, *ListCellRequest) (*ListCellReply, error)
	mustEmbedUnimplementedCellServer()
}

// UnimplementedCellServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCellServer struct{}

func (UnimplementedCellServer) CreateCell(context.Context, *CreateCellRequest) (*CreateCellReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCell not implemented")
}
func (UnimplementedCellServer) UpdateCell(context.Context, *UpdateCellRequest) (*UpdateCellReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCell not implemented")
}
func (UnimplementedCellServer) DeleteCell(context.Context, *DeleteCellRequest) (*DeleteCellReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCell not implemented")
}
func (UnimplementedCellServer) GetCell(context.Context, *GetCellRequest) (*GetCellReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCell not implemented")
}
func (UnimplementedCellServer) ListCell(context.Context, *ListCellRequest) (*ListCellReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCell not implemented")
}
func (UnimplementedCellServer) mustEmbedUnimplementedCellServer() {}
func (UnimplementedCellServer) testEmbeddedByValue()              {}

// UnsafeCellServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CellServer will
// result in compilation errors.
type UnsafeCellServer interface {
	mustEmbedUnimplementedCellServer()
}

func RegisterCellServer(s grpc.ServiceRegistrar, srv CellServer) {
	// If the following call pancis, it indicates UnimplementedCellServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Cell_ServiceDesc, srv)
}

func _Cell_CreateCell_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCellRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CellServer).CreateCell(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cell_CreateCell_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CellServer).CreateCell(ctx, req.(*CreateCellRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cell_UpdateCell_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCellRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CellServer).UpdateCell(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cell_UpdateCell_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CellServer).UpdateCell(ctx, req.(*UpdateCellRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cell_DeleteCell_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCellRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CellServer).DeleteCell(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cell_DeleteCell_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CellServer).DeleteCell(ctx, req.(*DeleteCellRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cell_GetCell_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCellRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CellServer).GetCell(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cell_GetCell_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CellServer).GetCell(ctx, req.(*GetCellRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cell_ListCell_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCellRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CellServer).ListCell(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cell_ListCell_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CellServer).ListCell(ctx, req.(*ListCellRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Cell_ServiceDesc is the grpc.ServiceDesc for Cell service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cell_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.helloworld.v1.Cell",
	HandlerType: (*CellServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCell",
			Handler:    _Cell_CreateCell_Handler,
		},
		{
			MethodName: "UpdateCell",
			Handler:    _Cell_UpdateCell_Handler,
		},
		{
			MethodName: "DeleteCell",
			Handler:    _Cell_DeleteCell_Handler,
		},
		{
			MethodName: "GetCell",
			Handler:    _Cell_GetCell_Handler,
		},
		{
			MethodName: "ListCell",
			Handler:    _Cell_ListCell_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "helloworld/v1/cell.proto",
}
