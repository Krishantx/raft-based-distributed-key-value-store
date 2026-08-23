package main

import (
	"fmt"
	"os"
	controller "raft-based-kv/internal/controller"
)

var memoryDB map[string]string

func init() {
	memoryDB = make(map[string]string)
}

type Config struct {
	NodeName string
}

func main() {
	go controller.StartHTTPController()
	var config Config = Config{}
	setConfig(&config)
	fmt.Println(config.NodeName)
}

func setConfig(config *Config) {
	config.NodeName = os.Getenv("node_name")
}
