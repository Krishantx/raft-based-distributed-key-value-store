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

// Send Heartbeat every 100ms

var node_name string
var term int
var conf models.Config

func StartLeaderService(config models.Config, channel chan bool) {
	term = 1
	node_name = config.NodeName
	conf = config
	ticker := time.NewTicker(3 * time.Second)
	for {
		<-ticker.C
		SendHeartbeat(channel)
	}
}

func grpcClient(addr string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewHeartbeatServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.LeaderDetails{NodeName: node_name}
	response, err := client.ReceiveHeartbeat(ctx, req)

	if err != nil {
		log.Fatalf("Error Sending Heartbear %s : %v", addr, err)
	}

	if int(response.Term) == term {
		fmt.Println("The current Term is correct and I am still the leader")
	} else {
		fmt.Println("The term does not match and I am no longer the leader")
	}
}

func SendHeartbeat(channel chan bool) {
	for i := 0; i < len(conf.Follower); i++ {
		addr := conf.Follower[i] + ":50051"
		grpcClient(addr)
	}
}
