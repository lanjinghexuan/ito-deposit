package server

import (
	"context"
	"fmt"
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

func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, order *service.OrderService, user *service.UserService,
	home *service.HomeService, deposit *service.DepositService, admin *service.AdminService, city *service.CityService, nearby *service.NearbyService,
	group *service.GroupService, cell *service.CabinetCellService, logger log.Logger) *http.Server {
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
	v1.RegisterGreeterHTTPServer(srv, greeter)
	v1.RegisterUserHTTPServer(srv, user)
	v1.RegisterHomeHTTPServer(srv, home)
	v1.RegisterDepositHTTPServer(srv, deposit)
	v1.RegisterOrderHTTPServer(srv, order)
	v1.RegisterCityHTTPServer(srv, city)
	v1.RegisterNearbyHTTPServer(srv, nearby)
	v1.RegisterAdminHTTPServer(srv, admin)
	v1.RegisterGroupHTTPServer(srv, group)
	v1.RegisterCabinetCellHTTPServer(srv, cell)

	if c.Pprof.Switch {
		go func() {
			fmt.Println(http2.ListenAndServe(fmt.Sprintf(":%d", c.Pprof.Port), nil))
		}()
	}

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
	whiteList["/api.helloworld.v1.User/GetUser"] = struct{}{}
	whiteList["/api.helloworld.v1.Order/ListOrder"] = struct{}{}
	whiteList["/api.helloworld.v1.Order/ShowOrder"] = struct{}{}
	whiteList["/api.helloworld.v1.Order/CreateOrder"] = struct{}{}
	whiteList["/api.helloworld.v1.Admin/PointList"] = struct{}{}
	whiteList["/api.helloworld.v1.Admin/PointInfo"] = struct{}{}
	whiteList["/api.helloworld.v1.Admin/SetPriceRule"] = struct{}{}
	whiteList["/api.helloworld.v1.Admin/GetPriceRule"] = struct{}{}
	whiteList["/api.helloworld.v1.Order/UpdateOrder"] = struct{}{}
	whiteList["/api.helloworld.v1.Admin/AdminLogin"] = struct{}{}
	// 寄存点相关API - 用户端不需要认证
	whiteList["/api.helloworld.v1.Deposit/GetDepositLocker"] = struct{}{}
	whiteList["/api.helloworld.v1.Deposit/CreateDeposit"] = struct{}{}
	whiteList["/api.helloworld.v1.Deposit/ListDeposit"] = struct{}{}
	// 城市相关API - 用户端不需要认证
	whiteList["/api.helloworld.v1.City/ListUserCities"] = struct{}{}
	whiteList["/api.helloworld.v1.City/SearchCities"] = struct{}{}
	whiteList["/api.helloworld.v1.City/GetUserCity"] = struct{}{}
	whiteList["/api.helloworld.v1.City/GetCityByCode"] = struct{}{}
	whiteList["/api.helloworld.v1.City/GetHotCities"] = struct{}{}
	// 附近寄存点相关API - 用户端不需要认证
	whiteList["/api.helloworld.v1.Nearby/FindNearbyLockerPoints"] = struct{}{}
	whiteList["/api.helloworld.v1.Nearby/FindNearbyLockerPointsInCity"] = struct{}{}
	whiteList["/api.helloworld.v1.Nearby/FindMyNearbyLockerPoints"] = struct{}{}
	whiteList["/api.helloworld.v1.Nearby/SearchLockerPointsInCity"] = struct{}{}
	whiteList["/api.helloworld.v1.Nearby/GetCityLockerPointsMap"] = struct{}{}
	whiteList["/api.helloworld.v1.Nearby/GetMyNearbyInfo"] = struct{}{}
	whiteList["/api.helloworld.v1.Nearby/GetAllLockerPoints"] = struct{}{}
	// 柜组相关API - 管理员功能，需要JWT验证，这里先添加到白名单用于测试
	whiteList["/api.helloworld.v1.Group/CreateGroup"] = struct{}{}
	whiteList["/api.helloworld.v1.Group/UpdateGroup"] = struct{}{}
	whiteList["/api.helloworld.v1.Group/DeleteGroup"] = struct{}{}
	whiteList["/api.helloworld.v1.Group/GetGroup"] = struct{}{}
	whiteList["/api.helloworld.v1.Group/ListGroup"] = struct{}{}
	whiteList["/api.helloworld.v1.Group/SearchGroup"] = struct{}{}
	// 柜格相关API - 管理员功能，需要JWT验证，这里先添加到白名单用于测试
	whiteList["/api.helloworld.v1.CabinetCell/CreateCabinetCell"] = struct{}{}
	whiteList["/api.helloworld.v1.CabinetCell/UpdateCabinetCell"] = struct{}{}
	whiteList["/api.helloworld.v1.CabinetCell/DeleteCabinetCell"] = struct{}{}
	whiteList["/api.helloworld.v1.CabinetCell/GetCabinetCell"] = struct{}{}
	whiteList["/api.helloworld.v1.CabinetCell/ListCabinetCells"] = struct{}{}
	whiteList["/api.helloworld.v1.CabinetCell/SearchCabinetCells"] = struct{}{}
	whiteList["/api.helloworld.v1.CabinetCell/GetCabinetCellsByGroup"] = struct{}{}
	whiteList["/api.helloworld.v1.CabinetCell/BatchCreateCabinetCells"] = struct{}{}
	whiteList["/api.helloworld.v1.CabinetCell/OpenCabinetCell"] = struct{}{}
	whiteList["/api.helloworld.v1.CabinetCell/CloseCabinetCell"] = struct{}{}

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
	loginWhiteList["/api.helloworld.v1.Admin/AdminLogin"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.Admin/CreateAdmin"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.Deposit/ReturnToken"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/SendSms"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/Register"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/Login"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/List"] = struct{}{}
	loginWhiteList["/api.helloworld.v1.User/Admin"] = struct{}{}

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
