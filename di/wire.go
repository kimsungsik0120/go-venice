//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-venice/configs"
	"go-venice/internal/adapter/nodes"
	"go-venice/internal/api/handler"
	"go-venice/internal/api/router"
	"go-venice/internal/api/service"
)

func InitializeRouter() *gin.Engine {
	wire.Build(
		configs.Load,
		nodes.NewEvm,
		wire.Bind(new(nodes.Node), new(*nodes.Evm)),
		service.NewNodeService,
		handler.NewNodeHandler,
		wire.Struct(new(router.RouterConfig), "*"),
		router.NewRouter,
	)
	return &gin.Engine{}
}
