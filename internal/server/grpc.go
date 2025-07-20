package server

import (
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	v1 "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, order *service.OrderService, user *service.UserService, home *service.HomeService, deposit *service.DepositService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	if false {
		opts = append(opts, grpc.Middleware(
			jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
				return []byte(c.Jwt.Authkey), nil
			}),
		))
	}

	srv := grpc.NewServer(opts...)
	v1.RegisterGreeterServer(srv, greeter)
	v1.RegisterUserServer(srv, user)
	v1.RegisterHomeServer(srv, home)
	v1.RegisterDepositServer(srv, deposit)
	v1.RegisterOrderServer(srv, order)
	return srv
}
