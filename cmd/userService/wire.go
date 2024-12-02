//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/TiktokCommence/userService/internal/biz"
	"github.com/TiktokCommence/userService/internal/conf"
	"github.com/TiktokCommence/userService/internal/data"
	"github.com/TiktokCommence/userService/internal/registry"
	"github.com/TiktokCommence/userService/internal/server"
	"github.com/TiktokCommence/userService/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.EmailConf, *conf.RegistryConf, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		service.ProviderSet,
		biz.ProviderSet,
		data.ProviderSet,
		registry.ProviderSet,
		newApp,
		wire.Bind(new(service.UserHandler), new(*biz.UserHandler)),
		wire.Bind(new(biz.GenerateID), new(*data.RedisWorkerImplement)),
		wire.Bind(new(biz.EmailWorker), new(*data.EmailWorker)),
		wire.Bind(new(biz.DBWorker), new(*data.UserRepo)),
		wire.Bind(new(biz.RedisWorker), new(*data.RedisWorkerImplement)),
	))
}
