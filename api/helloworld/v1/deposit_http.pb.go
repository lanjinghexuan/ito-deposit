// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             v5.26.1
// source: helloworld/v1/deposit.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationDepositCreateDeposit = "/api.helloworld.v1.Deposit/CreateDeposit"
const OperationDepositDecodeToken = "/api.helloworld.v1.Deposit/DecodeToken"
const OperationDepositGetDepositLocker = "/api.helloworld.v1.Deposit/GetDepositLocker"
const OperationDepositListDeposit = "/api.helloworld.v1.Deposit/ListDeposit"
const OperationDepositReturnToken = "/api.helloworld.v1.Deposit/ReturnToken"
const OperationDepositSendCodeByOrder = "/api.helloworld.v1.Deposit/SendCodeByOrder"
const OperationDepositUpdateDepositLockerId = "/api.helloworld.v1.Deposit/UpdateDepositLockerId"

type DepositHTTPServer interface {
	CreateDeposit(context.Context, *CreateDepositRequest) (*CreateDepositReply, error)
	DecodeToken(context.Context, *ReturnTokenReq) (*ReturnTokenRes, error)
	GetDepositLocker(context.Context, *GetDepositLockerReq) (*GetDepositLockerRes, error)
	ListDeposit(context.Context, *ListDepositRequest) (*ListDepositReply, error)
	ReturnToken(context.Context, *ReturnTokenReq) (*ReturnTokenRes, error)
	SendCodeByOrder(context.Context, *SendCodeByOrderReq) (*SendCodeByOrderRes, error)
	UpdateDepositLockerId(context.Context, *UpdateDepositLockerIdReq) (*UpdateDepositLockerIdRes, error)
}

func RegisterDepositHTTPServer(s *http.Server, srv DepositHTTPServer) {
	r := s.Route("/")
	r.POST("/deposit/createDeposit", _Deposit_CreateDeposit0_HTTP_Handler(srv))
	r.GET("/deposit", _Deposit_ListDeposit0_HTTP_Handler(srv))
	r.GET("/returntoken", _Deposit_ReturnToken0_HTTP_Handler(srv))
	r.GET("/decodetoken", _Deposit_DecodeToken0_HTTP_Handler(srv))
	r.GET("/getDepositLocker", _Deposit_GetDepositLocker0_HTTP_Handler(srv))
	r.POST("/deposit/updateDepositLockerId", _Deposit_UpdateDepositLockerId0_HTTP_Handler(srv))
	r.POST("/deposit/sendCodeByOrder", _Deposit_SendCodeByOrder0_HTTP_Handler(srv))
}

func _Deposit_CreateDeposit0_HTTP_Handler(srv DepositHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateDepositRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDepositCreateDeposit)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateDeposit(ctx, req.(*CreateDepositRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateDepositReply)
		return ctx.Result(200, reply)
	}
}

func _Deposit_ListDeposit0_HTTP_Handler(srv DepositHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListDepositRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDepositListDeposit)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListDeposit(ctx, req.(*ListDepositRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListDepositReply)
		return ctx.Result(200, reply)
	}
}

func _Deposit_ReturnToken0_HTTP_Handler(srv DepositHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ReturnTokenReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDepositReturnToken)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ReturnToken(ctx, req.(*ReturnTokenReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ReturnTokenRes)
		return ctx.Result(200, reply)
	}
}

func _Deposit_DecodeToken0_HTTP_Handler(srv DepositHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ReturnTokenReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDepositDecodeToken)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DecodeToken(ctx, req.(*ReturnTokenReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ReturnTokenRes)
		return ctx.Result(200, reply)
	}
}

func _Deposit_GetDepositLocker0_HTTP_Handler(srv DepositHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetDepositLockerReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDepositGetDepositLocker)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetDepositLocker(ctx, req.(*GetDepositLockerReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetDepositLockerRes)
		return ctx.Result(200, reply)
	}
}

func _Deposit_UpdateDepositLockerId0_HTTP_Handler(srv DepositHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateDepositLockerIdReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDepositUpdateDepositLockerId)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateDepositLockerId(ctx, req.(*UpdateDepositLockerIdReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateDepositLockerIdRes)
		return ctx.Result(200, reply)
	}
}

func _Deposit_SendCodeByOrder0_HTTP_Handler(srv DepositHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SendCodeByOrderReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDepositSendCodeByOrder)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SendCodeByOrder(ctx, req.(*SendCodeByOrderReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SendCodeByOrderRes)
		return ctx.Result(200, reply)
	}
}

type DepositHTTPClient interface {
	CreateDeposit(ctx context.Context, req *CreateDepositRequest, opts ...http.CallOption) (rsp *CreateDepositReply, err error)
	DecodeToken(ctx context.Context, req *ReturnTokenReq, opts ...http.CallOption) (rsp *ReturnTokenRes, err error)
	GetDepositLocker(ctx context.Context, req *GetDepositLockerReq, opts ...http.CallOption) (rsp *GetDepositLockerRes, err error)
	ListDeposit(ctx context.Context, req *ListDepositRequest, opts ...http.CallOption) (rsp *ListDepositReply, err error)
	ReturnToken(ctx context.Context, req *ReturnTokenReq, opts ...http.CallOption) (rsp *ReturnTokenRes, err error)
	SendCodeByOrder(ctx context.Context, req *SendCodeByOrderReq, opts ...http.CallOption) (rsp *SendCodeByOrderRes, err error)
	UpdateDepositLockerId(ctx context.Context, req *UpdateDepositLockerIdReq, opts ...http.CallOption) (rsp *UpdateDepositLockerIdRes, err error)
}

type DepositHTTPClientImpl struct {
	cc *http.Client
}

func NewDepositHTTPClient(client *http.Client) DepositHTTPClient {
	return &DepositHTTPClientImpl{client}
}

func (c *DepositHTTPClientImpl) CreateDeposit(ctx context.Context, in *CreateDepositRequest, opts ...http.CallOption) (*CreateDepositReply, error) {
	var out CreateDepositReply
	pattern := "/deposit/createDeposit"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationDepositCreateDeposit))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DepositHTTPClientImpl) DecodeToken(ctx context.Context, in *ReturnTokenReq, opts ...http.CallOption) (*ReturnTokenRes, error) {
	var out ReturnTokenRes
	pattern := "/decodetoken"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationDepositDecodeToken))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DepositHTTPClientImpl) GetDepositLocker(ctx context.Context, in *GetDepositLockerReq, opts ...http.CallOption) (*GetDepositLockerRes, error) {
	var out GetDepositLockerRes
	pattern := "/getDepositLocker"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationDepositGetDepositLocker))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DepositHTTPClientImpl) ListDeposit(ctx context.Context, in *ListDepositRequest, opts ...http.CallOption) (*ListDepositReply, error) {
	var out ListDepositReply
	pattern := "/deposit"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationDepositListDeposit))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DepositHTTPClientImpl) ReturnToken(ctx context.Context, in *ReturnTokenReq, opts ...http.CallOption) (*ReturnTokenRes, error) {
	var out ReturnTokenRes
	pattern := "/returntoken"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationDepositReturnToken))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DepositHTTPClientImpl) SendCodeByOrder(ctx context.Context, in *SendCodeByOrderReq, opts ...http.CallOption) (*SendCodeByOrderRes, error) {
	var out SendCodeByOrderRes
	pattern := "/deposit/sendCodeByOrder"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationDepositSendCodeByOrder))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DepositHTTPClientImpl) UpdateDepositLockerId(ctx context.Context, in *UpdateDepositLockerIdReq, opts ...http.CallOption) (*UpdateDepositLockerIdRes, error) {
	var out UpdateDepositLockerIdRes
	pattern := "/deposit/updateDepositLockerId"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationDepositUpdateDepositLockerId))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
