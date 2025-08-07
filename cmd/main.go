package main

import (
	"go-venice/configs"
	"go-venice/di"
)

func main() {
	// configs, logger, router 초기화
	config := configs.Load()

	//wire를 통한 주입
	router := di.InitializeRouter()

	//수동 의존성 주입
	/*
		node := nodes.NewEvm(config)
		nodeService := service.NewNodeService(config, node)
		router := router.NewRouter(router.RouterConfig{
			NodeService: nodeService,
		})
	*/
	if err := router.Run(":" + config.Port); err != nil {
		panic(err)
	}
}
