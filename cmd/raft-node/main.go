package main

import (
	"log"
	"raft-based-kv/internal/enums"
	"raft-based-kv/internal/loggers"
	models "raft-based-kv/internal/models"
	"raft-based-kv/internal/repo"
	service "raft-based-kv/internal/service"
)

func main() {
	// Use the Startup service to first
	// create and use all the vars and
	// struct. Then once that is d
	log.Println("Creating the repository")
	repo := repo.NewRepository()
	log.Println("Starting the bootup service to preload the repository if required")
	service.NewStartUpService().ModifyCurrentStateToLatestLog(*repo)
	TimerOps := make(chan enums.TimerOperations)
	config := models.Config{}
	config.Init()
	log.Println("config role : " + config.Role)
	nodeRaftData := models.NodeRaftData{}
	nodeRaftData.Init()
	loggers.Init()
	for {
		if config.Role == enums.Follower {
			log.Println("Role : Follower")
			follower_service := service.NewFollowerService(&config, TimerOps, &nodeRaftData)
			follower_service.StartFollowerService()
			follower_service.Close()
		} else if config.Role == enums.Leader {
			leader_service := service.NewLeaderService(config, nodeRaftData, *repo)
			log.Println("Role : Leader")
			leader_service.StartLeaderService()
			leader_service.Close()
		} else if config.Role == enums.Candidate {
			log.Println("Role : Candidate")
			candidateService := service.NewCandidateService(&config)
			candidateService.StartCandidateService(&nodeRaftData)
		}
	}
}
