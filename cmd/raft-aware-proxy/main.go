package main

import (
	controller "raft-based-kv/cmd/raft-aware-proxy/controller"
	service "raft-based-kv/cmd/raft-aware-proxy/service"
)

func main() {
	proxyService := service.NewProxyService()
	httpController := controller.NewHTTPController(proxyService)
	go httpController.StartHTTPController()

	select {}
}
