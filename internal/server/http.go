package server

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	v1 "ito-deposit/api/helloworld/v1"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/service"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, order *service.OrderService, user *service.UserService, home *service.HomeService, deposit *service.DepositService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	if true {
		opts = append(opts, http.Middleware(
			selector.Server(
				jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
					return []byte(c.Jwt.Authkey), nil
				}, jwt.WithSigningMethod(jwtv5.SigningMethodHS256), jwt.WithClaims(func() jwtv5.Claims {
					return &jwtv5.MapClaims{}
				})),
			).
				Match(NewWhiteListMatcher()).
				Build(),
		))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterUserHTTPServer(srv, user)
	v1.RegisterHomeHTTPServer(srv, home)
	v1.RegisterDepositHTTPServer(srv, deposit)
	v1.RegisterOrderHTTPServer(srv, order)
	return srv
}

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	// 添加不需要 JWT 验证的接口到白名单
	whiteList["/api.helloworld.v1.Deposit/ReturnToken"] = struct{}{}
	whiteList["/shop.interface.v1.ShopInterface/Register"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
