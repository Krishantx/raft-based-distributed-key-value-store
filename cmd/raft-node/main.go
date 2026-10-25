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
	TimerOps := make(chan enums.TimerOperations)
	config := models.Config{}
	config.Init()
	log.Println("config role : " + config.Role)
	nodeRaftData := models.NodeRaftData{}
	nodeRaftData.Init()
	loggers.Init()
	repo := repo.NewRepository()
	for {
		if config.Role == enums.Follower {
			log.Println("Role : Follower")
			follower_service := service.New_Follower_Service(&config, TimerOps)
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
