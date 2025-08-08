package router

import (
	"github.com/gin-gonic/gin"
	"go-venice/internal/api/handler"
)

func nodeRouter(router *gin.Engine, cfg RouterConfig) {
	nodeHandler := handler.NewNodeHandler(cfg.NodeService)

	v1 := router.Group("/v1")
	node := v1.Group("/node")
	{
		node.GET("/balance", nodeHandler.GetBalance)
	}
}
