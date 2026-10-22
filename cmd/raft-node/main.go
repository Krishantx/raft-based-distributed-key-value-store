package main

import (
	"fmt"
	"os"
	controller "raft-based-kv/internal/controller"
	"raft-based-kv/internal/enums"
	models "raft-based-kv/internal/models"
	repository "raft-based-kv/internal/repo"
	service "raft-based-kv/internal/service"
	"strconv"
	"strings"
)

func main() {
	TimerOps := make(chan enums.TimerOperations)

	config := setConfig()

	nodeRaftData := models.NodeRaftData{}

	nodeRaftData.Init()

	Repo := repository.NewRepository()

	grpc_Server := controller.New_gRPC_Server(Repo, &nodeRaftData, TimerOps)
	go grpc_Server.StartgRpcController(config)

	for {
		// Based on the role of the Node we have specific behaviours
		if config.Role == enums.Follower {
			// If the role of the Node is Follower. Start the follower service
			follower_service := service.New_Follower_Service(20, &config, TimerOps)
			follower_service.StartFollowerService()
		} else if config.Role == enums.Leader {
			//If the role of the Node is Leader. Start the Leader Service
			leader_service := service.NewLeaderService(config)
			leader_service.StartLeaderService()
		} else if config.Role == enums.Candidate {
			// If the role of the Node is Candidate. Start the follower service
			candidateService := service.NewCandidateService(&config)
			candidateService.StartCandidateService(&nodeRaftData)
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
	config.Role = enums.Role(os.Getenv("role"))
	config.Follower = strings.Split(os.Getenv("followers"), ",")
	return config
}
