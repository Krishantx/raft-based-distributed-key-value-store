package main

import (
	"fmt"
	"os"
	controller "raft-based-kv/internal/controller"
	models "raft-based-kv/internal/models"
	repository "raft-based-kv/internal/repo"
	"raft-based-kv/internal/service"
	"strconv"
	"strings"
)

var Term = 1
var Log = 1

func main() {
	config := setConfig()
	Repo := repository.NewRepository()
	grpc_Server := controller.New_gRPC_Server(Repo)
	go grpc_Server.StartgRpcController(config)

	for {
		if config.Role == "follower" {
			follower_service := service.New_Follower_Service(20)
			follower_service.StartFollowerService(config)
		} else {
			leader_service := service.NewLeaderService(config)
			leader_service.StartLeaderService()
		}
	}
}

func setConfig() models.Config {
	config := models.Config{}

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
	return config
}
