//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

func app() string {
	wire.Build(provider, ProvideNames)
	return ""
}
