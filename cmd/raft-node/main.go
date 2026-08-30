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

var Term = 1
var Log = 1

func main() {
	voteChan := make(chan bool, 1)
	channel := make(chan bool, 2)
	var config = models.Config{}
	setConfig(&config)

	grpc_Server := controller.New_gRPC_Server()
	go grpc_Server.StartgRpcController(config, channel, voteChan)

	for {
		if config.Role == "follower" {
			follower_service := service.New_Follower_Service(20)
			follower_service.StartFollowerService(config, channel, voteChan)
		} else {
			service.StartLeaderService(config, channel)
		}
	}
}

func setConfig(config *models.Config) {
	config.NodeName = os.Getenv("node_name")
	t := os.Getenv("election_time")
	fmt.Println("ElectionTime :" + t)
	time, err := strconv.Atoi(os.Getenv("election_time"))
	if err != nil {
		fmt.Printf("Cannot parse Election Time: %s", err)
	}
	config.ElectionTime = time
	config.Role = os.Getenv("role")
	config.Follower = strings.Split(os.Getenv("followers"), ",")
}
