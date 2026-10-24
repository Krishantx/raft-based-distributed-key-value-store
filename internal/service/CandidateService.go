package service

import (
	"context"
	"fmt"
	"log"
	"raft-based-kv/internal/clients"
	"raft-based-kv/internal/enums"
	"raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
	"time"

	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
)

type CandidateService struct {
	grpcClients clients.GRPCClientList
	config      *models.Config
}

func NewCandidateService(config *models.Config) CandidateService {
	return CandidateService{
		config: config,
	}
}

func (this *CandidateService) StartCandidateService(nodeRaftData *models.NodeRaftData) {
	// Send gRPC request to every single node.
	//  For this I also want to know what Term I am on and what index I am on.

	req := &pb.NodeInfo{
		Term: int32(nodeRaftData.GetTerm()),
		Log:  int32(nodeRaftData.GetIndex()),
	}

	voteResponse := [5]pb.Vote{}

	//Gather Votes
	for i := 0; i < 5; i++ {
		nodeRaftData.IncrementTerm()
		addr := "raft-node-1" + ":" + "50051"
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect : %s: %v", addr, err)
		}

		client := pb.NewVotingServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		res, err := client.StartVoting(ctx, req)
		if err != nil {
			log.Panic("Error Sending gRPC Request", err)
		}
		voteResponse[i] = *res
	}
	// Count Votes and Become Leader if majority voted for you
	voteCount := 0
	for i := 0; i < 5; i++ {
		if voteResponse[i].Vote == true {
			voteCount++
		}
	}

	if voteCount > 2 {
		// Become Leader
		fmt.Println("Became Leader")
		this.config.Role = enums.Leader
	} else {
		// Become Follower
		fmt.Println("Become Follower")
		this.config.Role = enums.Follower
	}

}
