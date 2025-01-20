//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"kraots-xa/configs"
	"kraots-xa/internal/conf"
	controller "kraots-xa/internal/control"
)

func WireApp(allConfig *configs.AllConfig) *conf.KratosServer {
	wire.Build(
		wire.Struct(new(conf.KratosServer), "*"),
		conf.ProviderSet,
		controller.ControlProvider,
	)
	return &conf.KratosServer{}
}
