// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.26.1
// source: helloworld/v1/deposit.proto

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
	Deposit_CreateDeposit_FullMethodName         = "/api.helloworld.v1.Deposit/CreateDeposit"
	Deposit_UpdateDeposit_FullMethodName         = "/api.helloworld.v1.Deposit/UpdateDeposit"
	Deposit_DeleteDeposit_FullMethodName         = "/api.helloworld.v1.Deposit/DeleteDeposit"
	Deposit_GetDeposit_FullMethodName            = "/api.helloworld.v1.Deposit/GetDeposit"
	Deposit_ListDeposit_FullMethodName           = "/api.helloworld.v1.Deposit/ListDeposit"
	Deposit_ReturnToken_FullMethodName           = "/api.helloworld.v1.Deposit/ReturnToken"
	Deposit_DecodeToken_FullMethodName           = "/api.helloworld.v1.Deposit/DecodeToken"
	Deposit_GetDepositLocker_FullMethodName      = "/api.helloworld.v1.Deposit/GetDepositLocker"
	Deposit_UpdateDepositLockerId_FullMethodName = "/api.helloworld.v1.Deposit/UpdateDepositLockerId"
	Deposit_SendCodeByOrder_FullMethodName       = "/api.helloworld.v1.Deposit/SendCodeByOrder"
)

// DepositClient is the client API for Deposit service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DepositClient interface {
	CreateDeposit(ctx context.Context, in *CreateDepositRequest, opts ...grpc.CallOption) (*CreateDepositReply, error)
	UpdateDeposit(ctx context.Context, in *UpdateDepositRequest, opts ...grpc.CallOption) (*UpdateDepositReply, error)
	DeleteDeposit(ctx context.Context, in *DeleteDepositRequest, opts ...grpc.CallOption) (*DeleteDepositReply, error)
	GetDeposit(ctx context.Context, in *GetDepositRequest, opts ...grpc.CallOption) (*GetDepositReply, error)
	ListDeposit(ctx context.Context, in *ListDepositRequest, opts ...grpc.CallOption) (*ListDepositReply, error)
	ReturnToken(ctx context.Context, in *ReturnTokenReq, opts ...grpc.CallOption) (*ReturnTokenRes, error)
	DecodeToken(ctx context.Context, in *ReturnTokenReq, opts ...grpc.CallOption) (*ReturnTokenRes, error)
	GetDepositLocker(ctx context.Context, in *GetDepositLockerReq, opts ...grpc.CallOption) (*GetDepositLockerRes, error)
	UpdateDepositLockerId(ctx context.Context, in *UpdateDepositLockerIdReq, opts ...grpc.CallOption) (*UpdateDepositLockerIdRes, error)
	SendCodeByOrder(ctx context.Context, in *SendCodeByOrderReq, opts ...grpc.CallOption) (*SendCodeByOrderRes, error)
}

type depositClient struct {
	cc grpc.ClientConnInterface
}

func NewDepositClient(cc grpc.ClientConnInterface) DepositClient {
	return &depositClient{cc}
}

func (c *depositClient) CreateDeposit(ctx context.Context, in *CreateDepositRequest, opts ...grpc.CallOption) (*CreateDepositReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateDepositReply)
	err := c.cc.Invoke(ctx, Deposit_CreateDeposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depositClient) UpdateDeposit(ctx context.Context, in *UpdateDepositRequest, opts ...grpc.CallOption) (*UpdateDepositReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateDepositReply)
	err := c.cc.Invoke(ctx, Deposit_UpdateDeposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depositClient) DeleteDeposit(ctx context.Context, in *DeleteDepositRequest, opts ...grpc.CallOption) (*DeleteDepositReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteDepositReply)
	err := c.cc.Invoke(ctx, Deposit_DeleteDeposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depositClient) GetDeposit(ctx context.Context, in *GetDepositRequest, opts ...grpc.CallOption) (*GetDepositReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDepositReply)
	err := c.cc.Invoke(ctx, Deposit_GetDeposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depositClient) ListDeposit(ctx context.Context, in *ListDepositRequest, opts ...grpc.CallOption) (*ListDepositReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListDepositReply)
	err := c.cc.Invoke(ctx, Deposit_ListDeposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depositClient) ReturnToken(ctx context.Context, in *ReturnTokenReq, opts ...grpc.CallOption) (*ReturnTokenRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReturnTokenRes)
	err := c.cc.Invoke(ctx, Deposit_ReturnToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depositClient) DecodeToken(ctx context.Context, in *ReturnTokenReq, opts ...grpc.CallOption) (*ReturnTokenRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReturnTokenRes)
	err := c.cc.Invoke(ctx, Deposit_DecodeToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depositClient) GetDepositLocker(ctx context.Context, in *GetDepositLockerReq, opts ...grpc.CallOption) (*GetDepositLockerRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDepositLockerRes)
	err := c.cc.Invoke(ctx, Deposit_GetDepositLocker_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depositClient) UpdateDepositLockerId(ctx context.Context, in *UpdateDepositLockerIdReq, opts ...grpc.CallOption) (*UpdateDepositLockerIdRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateDepositLockerIdRes)
	err := c.cc.Invoke(ctx, Deposit_UpdateDepositLockerId_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *depositClient) SendCodeByOrder(ctx context.Context, in *SendCodeByOrderReq, opts ...grpc.CallOption) (*SendCodeByOrderRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendCodeByOrderRes)
	err := c.cc.Invoke(ctx, Deposit_SendCodeByOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DepositServer is the server API for Deposit service.
// All implementations must embed UnimplementedDepositServer
// for forward compatibility.
type DepositServer interface {
	CreateDeposit(context.Context, *CreateDepositRequest) (*CreateDepositReply, error)
	UpdateDeposit(context.Context, *UpdateDepositRequest) (*UpdateDepositReply, error)
	DeleteDeposit(context.Context, *DeleteDepositRequest) (*DeleteDepositReply, error)
	GetDeposit(context.Context, *GetDepositRequest) (*GetDepositReply, error)
	ListDeposit(context.Context, *ListDepositRequest) (*ListDepositReply, error)
	ReturnToken(context.Context, *ReturnTokenReq) (*ReturnTokenRes, error)
	DecodeToken(context.Context, *ReturnTokenReq) (*ReturnTokenRes, error)
	GetDepositLocker(context.Context, *GetDepositLockerReq) (*GetDepositLockerRes, error)
	UpdateDepositLockerId(context.Context, *UpdateDepositLockerIdReq) (*UpdateDepositLockerIdRes, error)
	SendCodeByOrder(context.Context, *SendCodeByOrderReq) (*SendCodeByOrderRes, error)
	mustEmbedUnimplementedDepositServer()
}

// UnimplementedDepositServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDepositServer struct{}

func (UnimplementedDepositServer) CreateDeposit(context.Context, *CreateDepositRequest) (*CreateDepositReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDeposit not implemented")
}
func (UnimplementedDepositServer) UpdateDeposit(context.Context, *UpdateDepositRequest) (*UpdateDepositReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDeposit not implemented")
}
func (UnimplementedDepositServer) DeleteDeposit(context.Context, *DeleteDepositRequest) (*DeleteDepositReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDeposit not implemented")
}
func (UnimplementedDepositServer) GetDeposit(context.Context, *GetDepositRequest) (*GetDepositReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeposit not implemented")
}
func (UnimplementedDepositServer) ListDeposit(context.Context, *ListDepositRequest) (*ListDepositReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDeposit not implemented")
}
func (UnimplementedDepositServer) ReturnToken(context.Context, *ReturnTokenReq) (*ReturnTokenRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReturnToken not implemented")
}
func (UnimplementedDepositServer) DecodeToken(context.Context, *ReturnTokenReq) (*ReturnTokenRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DecodeToken not implemented")
}
func (UnimplementedDepositServer) GetDepositLocker(context.Context, *GetDepositLockerReq) (*GetDepositLockerRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDepositLocker not implemented")
}
func (UnimplementedDepositServer) UpdateDepositLockerId(context.Context, *UpdateDepositLockerIdReq) (*UpdateDepositLockerIdRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDepositLockerId not implemented")
}
func (UnimplementedDepositServer) SendCodeByOrder(context.Context, *SendCodeByOrderReq) (*SendCodeByOrderRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCodeByOrder not implemented")
}
func (UnimplementedDepositServer) mustEmbedUnimplementedDepositServer() {}
func (UnimplementedDepositServer) testEmbeddedByValue()                 {}

// UnsafeDepositServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DepositServer will
// result in compilation errors.
type UnsafeDepositServer interface {
	mustEmbedUnimplementedDepositServer()
}

func RegisterDepositServer(s grpc.ServiceRegistrar, srv DepositServer) {
	// If the following call pancis, it indicates UnimplementedDepositServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Deposit_ServiceDesc, srv)
}

func _Deposit_CreateDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).CreateDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_CreateDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).CreateDeposit(ctx, req.(*CreateDepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deposit_UpdateDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).UpdateDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_UpdateDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).UpdateDeposit(ctx, req.(*UpdateDepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deposit_DeleteDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).DeleteDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_DeleteDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).DeleteDeposit(ctx, req.(*DeleteDepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deposit_GetDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).GetDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_GetDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).GetDeposit(ctx, req.(*GetDepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deposit_ListDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).ListDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_ListDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).ListDeposit(ctx, req.(*ListDepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deposit_ReturnToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReturnTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).ReturnToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_ReturnToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).ReturnToken(ctx, req.(*ReturnTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deposit_DecodeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReturnTokenReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).DecodeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_DecodeToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).DecodeToken(ctx, req.(*ReturnTokenReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deposit_GetDepositLocker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDepositLockerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).GetDepositLocker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_GetDepositLocker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).GetDepositLocker(ctx, req.(*GetDepositLockerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deposit_UpdateDepositLockerId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDepositLockerIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).UpdateDepositLockerId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_UpdateDepositLockerId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).UpdateDepositLockerId(ctx, req.(*UpdateDepositLockerIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deposit_SendCodeByOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendCodeByOrderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepositServer).SendCodeByOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Deposit_SendCodeByOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepositServer).SendCodeByOrder(ctx, req.(*SendCodeByOrderReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Deposit_ServiceDesc is the grpc.ServiceDesc for Deposit service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Deposit_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.helloworld.v1.Deposit",
	HandlerType: (*DepositServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateDeposit",
			Handler:    _Deposit_CreateDeposit_Handler,
		},
		{
			MethodName: "UpdateDeposit",
			Handler:    _Deposit_UpdateDeposit_Handler,
		},
		{
			MethodName: "DeleteDeposit",
			Handler:    _Deposit_DeleteDeposit_Handler,
		},
		{
			MethodName: "GetDeposit",
			Handler:    _Deposit_GetDeposit_Handler,
		},
		{
			MethodName: "ListDeposit",
			Handler:    _Deposit_ListDeposit_Handler,
		},
		{
			MethodName: "ReturnToken",
			Handler:    _Deposit_ReturnToken_Handler,
		},
		{
			MethodName: "DecodeToken",
			Handler:    _Deposit_DecodeToken_Handler,
		},
		{
			MethodName: "GetDepositLocker",
			Handler:    _Deposit_GetDepositLocker_Handler,
		},
		{
			MethodName: "UpdateDepositLockerId",
			Handler:    _Deposit_UpdateDepositLockerId_Handler,
		},
		{
			MethodName: "SendCodeByOrder",
			Handler:    _Deposit_SendCodeByOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "helloworld/v1/deposit.proto",
}
