package main

import (
	"fmt"
	"os"
	controller "raft-based-kv/internal/controller"
	service "raft-based-kv/internal/service"
	"strconv"
)

// Change Role type to enum later
type Config struct {
	NodeName     string
	ElectionTime int
	Role         string
}

func main() {
	var config = Config{}
	setConfig(&config)

	go controller.StartgRpcController()
	go service.StartElectionTimer(config.ElectionTime)
}

func setConfig(config *Config) {
	config.NodeName = os.Getenv("node_name")
	time, err := strconv.Atoi(os.Getenv("election_time"))
	if err != nil {
		fmt.Printf("Cannot parse Election Time: %s", err)
	}
	config.ElectionTime = time
	config.Role = os.Getenv("role")

}
