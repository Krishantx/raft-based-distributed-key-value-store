package service

import (
	"context"
	"fmt"
	"log"
	"raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CandidateService struct {
}

func NewCandidateService() CandidateService {
	return CandidateService{}
}

func (this *CandidateService) StartCandidateService(nodeRaftData *models.NodeRaftData) {
	// Send gRPC request to every single node.
	//  For this I also want to know what Term I am on and what index I am on.

	req := &pb.NodeInfo{
		Term: int32(nodeRaftData.GetTerm()),
		Log:  int32(nodeRaftData.GetIndex()),
	}
	client, ctx, _ := createGrpcClinet("raft-node-1:80081")
	response, err := client.StartVoting(ctx, req)

	if err != nil {
		log.Fatalf("Error Sending Heartbear %s : %v", "raft-node-1", err)
	}
	fmt.Print(response)
}

func createGrpcClinet(addr string) (pb.VotingServiceClient, context.Context, context.CancelFunc) {

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewVotingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return client, ctx, cancel
}
