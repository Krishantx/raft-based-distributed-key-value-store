package main

import (
	controller "raft-based-kv/internal/controller"
	"raft-based-kv/internal/models"
	"raft-based-kv/internal/service"
)

type context struct {
	leader models.Node
}

func main() {
	proxyService := service.NewProxyService()
	httpController := controller.NewHTTPController(proxyService)
	go httpController.StartHTTPController(proxyService)

	select {}
}
