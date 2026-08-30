package main

import (
	"fmt"
	"os"
	controller "raft-based-kv/internal/controller"
	models "raft-based-kv/internal/models"
	"raft-based-kv/internal/service"
	"strconv"
	"strings"
)

var term = 1
var log = 1

func main() {
	channel := make(chan bool, 2)
	var config = models.Config{}
	setConfig(&config)

	go controller.StartgRpcController(config, channel)
	for {
		if config.Role == "follower" {
			service.StartFollowerService(config, channel)
		} else {
			service.StartLeaderService(config, channel)
		}
	}
}

func setConfig(config *models.Config) {
	config.NodeName = os.Getenv("node_name")
	time, err := strconv.Atoi(os.Getenv("election_time"))
	if err != nil {
		fmt.Printf("Cannot parse Election Time: %s", err)
	}
	config.ElectionTime = time
	config.Role = os.Getenv("role")
	config.Follower = strings.Split(os.Getenv("followers"), ",")
}
