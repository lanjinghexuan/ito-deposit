//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"ito-deposit/internal/biz"
	"ito-deposit/internal/conf"
	"ito-deposit/internal/data"
	"ito-deposit/internal/server"
	"ito-deposit/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"ito-deposit/internal/pkg/geo"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp, geo.ProviderSet))
}
