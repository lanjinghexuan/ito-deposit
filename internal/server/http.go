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
	http2 "net/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, order *service.OrderService, user *service.UserService, home *service.HomeService, deposit *service.DepositService, city *service.CityService, nearby *service.NearbyService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Filter(corsFilter),
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
	v1.RegisterUserHTTPServer(srv, user)
	v1.RegisterHomeHTTPServer(srv, home)
	v1.RegisterDepositHTTPServer(srv, deposit)
	v1.RegisterOrderHTTPServer(srv, order)
	v1.RegisterCityHTTPServer(srv, city)
	v1.RegisterNearbyHTTPServer(srv, nearby)
	return srv
}

func NewWhiteListMatcher() selector.MatchFunc {
	// 创建需要JWT验证的接口列表（黑名单）
	// 只有管理员相关的API需要JWT验证
	jwtRequiredList := make(map[string]struct{})

	// 管理员API需要JWT验证
	jwtRequiredList["/api.helloworld.v1.Admin/GetAdminInfo"] = struct{}{}
	jwtRequiredList["/api.helloworld.v1.Admin/UpdateAdmin"] = struct{}{}
	jwtRequiredList["/api.helloworld.v1.Admin/DeleteAdmin"] = struct{}{}
	jwtRequiredList["/api.helloworld.v1.Admin/ListAdmins"] = struct{}{}

	// 附近寄存点管理API需要JWT验证
	jwtRequiredList["/api.helloworld.v1.Nearby/InitLockerPointsGeo"] = struct{}{}

	// 其他需要JWT验证的管理API
	// ...

	// 特殊情况：管理员登录和创建管理员不需要JWT验证
	loginWhiteList := make(map[string]struct{})
	loginWhiteList["/api.helloworld.v1.Admin/Login"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.Admin/CreateAdmin"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.Deposit/ReturnToken"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/SendSms"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/Register"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/Login"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/List"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/Admin"] = struct{}{}

	return func(ctx context.Context, operation string) bool {
		// 检查是否是管理员登录或创建管理员API
		if _, ok := loginWhiteList[operation]; ok {
			return false // 不需要JWT验证
		}

		// 检查是否是管理员API
		if _, ok := jwtRequiredList[operation]; ok {
			return true // 需要JWT验证
		}

		// 默认情况下，所有其他API不需要JWT验证
		return false
	}
}

// corsFilter 手动实现 CORS
func corsFilter(next http2.Handler) http2.Handler {
	return http2.HandlerFunc(func(w http2.ResponseWriter, r *http2.Request) {
		// 允许的域名、方法、头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		// 预检请求直接返回 204
		if r.Method == http2.MethodOptions {
			w.WriteHeader(http2.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
