package main

import (
	"raft-based-kv/internal/controller"
)

func main() {
	go controller.StartHTTPController()

}
