//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"ito-deposit/internal/basic/pkg/job"
	"ito-deposit/internal/biz"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/data"
	"ito-deposit/internal/pkg/cronserver"
	"ito-deposit/internal/server"
	"ito-deposit/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
<<<<<<< HEAD
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

// newApp 创建应用实例
func newApp(gs *grpc.Server, hs *http.Server, cs *cronserver.Server, logger log.Logger) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
			cs, // 添加定时任务服务器
		),
	)
=======
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp, NewRegistrar, NewEtcdClient, geo.ProviderSet, NewContext, job.NewScheduler))
>>>>>>> c7faa8141686d333f091a98906bccc7ba10312da
}
