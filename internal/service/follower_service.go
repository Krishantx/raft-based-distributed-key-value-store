package service

import (
	"context"
	"fmt"
	"log"
	"raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
	"time"

	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
)

func New_Follower_Service(VotingTimeout int, config *models.Config) Follower_Service {

	return Follower_Service{
		VotingTimeout: VotingTimeout,
		config:        config,
	}
}

type Follower_Service struct {
	VotingTimeout int
	config        *models.Config
}

func (this *Follower_Service) StartFollowerService() {
	this.StartElectionTimer()
}

func (this *Follower_Service) StartElectionTimer() {
	fmt.Println("Election Timer Started")
	timer := time.NewTimer(5 * time.Second)
	for this.config.Role == models.Follower {
		select {
		//case 1: Leader did not send a Heartbeat, start voting.
		case <-timer.C:
			fmt.Println("Timer expired: Changing role to candidate")
			this.config.Role = models.Candidate
			break
			// case <-ch:
			// fmt.Println("Heartbeat Recieved Reset Timer")
			// timer.Reset(5 * time.Second)
		}
	}
}
func grpcClientForVoting(addr string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewVotingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.NodeInfo{
		Term: 1,
		Log:  1,
		// NodeName: conf.NodeName,
	}
	response, err := client.StartVoting(ctx, req)

	if err != nil {
		log.Fatalf("Error Sending Heartbear %s : %v", addr, err)
	}
	fmt.Print(response)
}
