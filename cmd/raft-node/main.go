package main

import (
	"fmt"
	"os"
	"raft-based-kv/internal/enums"
	"raft-based-kv/internal/loggers"
	models "raft-based-kv/internal/models"
	"raft-based-kv/internal/repo"
	service "raft-based-kv/internal/service"
	"strconv"
	"strings"
)

func main() {
	TimerOps := make(chan enums.TimerOperations)
	config := setConfig()
	nodeRaftData := models.NodeRaftData{}
	nodeRaftData.Init()
	loggers.Init()
	repo := repo.NewRepository()
	for {
		if config.Role == enums.Follower {
			follower_service := service.New_Follower_Service(20, &config, TimerOps)
			follower_service.StartFollowerService()
			follower_service.Close()
		} else if config.Role == enums.Leader {
			leader_service := service.NewLeaderService(config, nodeRaftData, *repo)
			leader_service.StartLeaderService()
			leader_service.Close()
		} else if config.Role == enums.Candidate {
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
