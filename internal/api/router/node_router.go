package router

import (
	"github.com/gin-gonic/gin"
)

func nodeRouter(router *gin.Engine, cfg RouterConfig) {
	v1 := router.Group("/v1")
	node := v1.Group("/node")
	{
		node.GET("/balance", cfg.NodeHandler.GetBalance)
	}
}
