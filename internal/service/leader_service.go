package service

import (
	"context"
	// "fmt"
	"log"
	"raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Send Heartbeat every 100ms

type LeaderService struct {
	node_name string
	term      int
	conf      models.Config
}

func NewLeaderService(config models.Config) LeaderService {
	return LeaderService{
		node_name: config.NodeName,
		term:      1,
		conf:      config,
	}
}

func (s *LeaderService) StartLeaderService() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		s.SendHeartbeat()
	}
}

func (this *LeaderService) grpcClient(addr string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewHeartbeatServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.Leader{NodeName: this.node_name}
	response, err := client.ReceiveHeartbeat(ctx, req)

	if err != nil {
		log.Fatalf("Error Sending Heartbear %s : %v", addr, err)
	}

	if int(response.Term) == this.term {
		// fmt.Println("The current Term is correct and I am still the leader")
	} else {
		// fmt.Println("The term does not match and I am no longer the leader")
	}
}

func (this *LeaderService) SendHeartbeat() {
	for i := 0; i < len(this.conf.Follower); i++ {
		addr := this.conf.Follower[i] + ":50051"
		this.grpcClient(addr)
	}
}
