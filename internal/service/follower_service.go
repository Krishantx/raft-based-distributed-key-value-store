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

func StartFollowerService(config models.Config, ch chan bool) {
	StartElectionTimer(config, ch)
}

func StartElectionTimer(config models.Config, ch chan bool) {
	timer := time.NewTimer(5 * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Println("Timer expired starting voting procedure")
			StartVoting()
		case <-ch:
			fmt.Println("Heartbeat Recieved Reset Timer")
			timer.Reset(5 * time.Second)
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
		Term:     1,
		Log:      1,
		NodeName: conf.NodeName,
	}
	response, err := client.StartVoting(ctx, req)

	if err != nil {
		log.Fatalf("Error Sending Heartbear %s : %v", addr, err)
	}
	fmt.Print(response)
}

func StartVoting() {
	nodes := [...]string{"raft-node-2", "raft-node-3", "raft-node-4", "raft-node-5"}
	for i := range len(nodes) {
		grpcClientForVoting(nodes[i] + ":50051")
	}
}
