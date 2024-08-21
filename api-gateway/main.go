package main

import (
	router "gitlab.com/bahodirova/api-gateway/api"
	"gitlab.com/bahodirova/api-gateway/api/handler"
	"gitlab.com/bahodirova/api-gateway/config"
)

func main() {

	engine := router.NewRouter(handler.NewHandlerStruct())
	engine.Run(config.Load().HTTPPort)
}
