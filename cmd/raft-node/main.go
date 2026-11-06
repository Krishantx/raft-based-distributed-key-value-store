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
	loggers.Init()
	log.Println("Creating the repository for in memory storage")
	repo := repo.NewRepository()
	log.Println("Starting the bootup service to preload the repository, if required, and fetching previous memory state")
	service.NewStartUpService(repo).ModifyCurrentStateToLatestLog()
	config := models.Config{}
	config.Init()
	log.Println("config role : " + config.Role)
	nodeRaftData := models.NodeRaftData{}
	nodeRaftData.Init()
	for {
		if config.Role == enums.Follower {
			log.Println("Role : Follower")
			follower_service := service.NewFollowerService(&config, &nodeRaftData)
			follower_service.StartFollowerService()
			follower_service.Close()
		} else if config.Role == enums.Leader {
			leader_service := service.NewLeaderService(&config, nodeRaftData, *repo)
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
