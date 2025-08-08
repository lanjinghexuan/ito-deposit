//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"ito-deposit/internal/basic/pkg/job"
	"ito-deposit/internal/biz"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/data"
	"ito-deposit/internal/pkg/geo"
	"ito-deposit/internal/server"
	"ito-deposit/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {

	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp, NewRegistrar, NewEtcdClient, geo.ProviderSet, NewContext, job.NewScheduler))

}
