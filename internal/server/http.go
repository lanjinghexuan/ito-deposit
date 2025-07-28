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
	_ "net/http/pprof"
)

func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, order *service.OrderService, user *service.UserService,
	home *service.HomeService, deposit *service.DepositService, admin *service.AdminService,
	logger log.Logger) *http.Server {
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

	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterUserHTTPServer(srv, user)
	v1.RegisterHomeHTTPServer(srv, home)
	v1.RegisterDepositHTTPServer(srv, deposit)
	v1.RegisterOrderHTTPServer(srv, order)
	v1.RegisterAdminHTTPServer(srv, admin)

	http2.ListenAndServe(":6001", nil)

	srv.Route("/").POST("/upload", admin.DownloadFile)

	return srv
}

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	// 添加不需要 JWT 验证的接口到白名单
	whiteList["/api.helloworld.v1.Deposit/ReturnToken"] = struct{}{}
	whiteList["/api.helloworld.v1.User/SendSms"] = struct{}{}
	whiteList["/api.helloworld.v1.User/Register"] = struct{}{}
	whiteList["/api.helloworld.v1.User/Login"] = struct{}{}
	whiteList["/api.helloworld.v1.User/List"] = struct{}{}
	whiteList["/api.helloworld.v1.User/Admin"] = struct{}{}
	whiteList["/api.helloworld.v1.Admin/AdminLogin"] = struct{}{}
	whiteList["/api.helloworld.v1.Order/ListOrder"] = struct{}{}
	whiteList["/api.helloworld.v1.Order/ShowOrder"] = struct{}{}
	whiteList["/api.helloworld.v1.Admin/PointList"] = struct{}{}
	whiteList["/api.helloworld.v1.Admin/PointInfo"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
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
