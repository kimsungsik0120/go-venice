package router

import (
	"github.com/gin-gonic/gin"
	"go-venice/internal/api/handler"
	"go-venice/internal/api/middleware"
)

type RouterConfig struct {
	NodeHandler handler.NodeHandler
	// 다른 서비스 주입 가능
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	engine := gin.New()

	// 미들웨어 등록
	engine.Use(middleware.Logging())
	engine.Use(middleware.Logging())
	engine.Use(gin.Recovery())

	// 라우트 등록
	nodeRouter(engine, cfg)

	return engine
}
