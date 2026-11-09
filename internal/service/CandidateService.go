package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"raft-based-kv/internal/clients"
	"raft-based-kv/internal/enums"
	"raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
	"strconv"
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

func (this *CandidateService) getVoteFromFollwer(nodeHostname string, req *pb.NodeInfo) (*pb.Vote, error) {
	addr := nodeHostname + ":" + "50051"
	if addr != this.config.NodeName+":50051" {
		res := &pb.Vote{}
		err := errors.New("")
		go func() {
			log.Printf("Sending a gRPC request to:  %v", addr)
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Failed to connect : %s: %v", addr, err)
			}
			client := pb.NewVotingServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			res, err = client.StartVoting(ctx, req)
		}()
		return res, err
	}
	return &pb.Vote{
		Vote: true,
	}, nil
}

func (this *CandidateService) StartCandidateService(nodeRaftData *models.NodeRaftData) {

	log.Println("=========Candidate Service Started=========")

	req := &pb.NodeInfo{
		Term: int32(nodeRaftData.GetTerm()),
		Log:  int32(nodeRaftData.GetIndex()),
	}

	voteResponse := [5]*pb.Vote{}

	nodeRaftData.IncrementTerm()
	//Gather Votes
	for i := 1; i <= 5; i++ {
		nodeHostname := "raft-node-" + strconv.Itoa(i)
		res, err := this.getVoteFromFollwer(nodeHostname, req)
		if err != nil {
			log.Printf("Error recieving vote from node %v : %v", i, err)
			voteResponse[i-1] = &pb.Vote{
				Vote: false,
			}
		} else {
			voteResponse[i-1] = res
		}

	}
	voteCount := this.countVotes(&voteResponse)
	this.config.Role = this.getRoleFromVoteCount(voteCount)
}

func (this *CandidateService) getRoleFromVoteCount(voteCount int) enums.Role {
	if voteCount > 2 {
		// Become Leader
		fmt.Println("Became Leader")
		return enums.Leader
	} else {
		fmt.Println("Become Follower")
		return enums.Follower
	}
}
func (this *CandidateService) countVotes(voteResponse *[5]*pb.Vote) int {
	voteCount := 0
	for i := 0; i < 5; i++ {
		if voteResponse[i].Vote == true {
			voteCount++
		}
	}
	return voteCount
}
