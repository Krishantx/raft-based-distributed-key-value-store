package main

import (
	controller "raft-based-kv/internal/controller"
	"raft-based-kv/internal/service"
)

func main() {
	proxyService := service.NewProxyService()
	httpController := controller.NewHTTPController(proxyService)
	go httpController.StartHTTPController(proxyService)

	select {}
}
