package service

import (
	"context"
	"fmt"
	"log"
	models "raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProxyService struct {
	Leader models.Node
}

func NewProxyService() ProxyService {
	node := models.Node{
		Name:     "A",
		Hostname: "raft-node-1",
		Port:     "50051",
	}
	return ProxyService{
		Leader: node,
	}
}

func (p *ProxyService) GetKeyValue() models.KeyValue {
	// Now use gRPC to send this get request to the Leader(Or who we belive is the leader)
	// First we fetch who we think the leader is
	leaderAddr := p.Leader.Hostname + ":" + p.Leader.Port
	// Send a gRPC Request to the Leader
	fmt.Println(leaderAddr)

	conn, err := grpc.NewClient(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", &leaderAddr, err)
	}

	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.Key{
		Term: 1,
		Key:  "key",
	}

	response, err := client.GetKeyValue(ctx, req)

	fmt.Println(response)

	return models.KeyValue{
		Key:   "This is a random Key",
		Value: "This is a random value",
	}
}

func (p *ProxyService) PutKeyValue() models.KeyValue {
	fmt.Printf("Hello World")
	return models.KeyValue{
		Key:   "This is a random Key",
		Value: "This is a random value",
	}
}
func (p *ProxyService) DeleteKeyValue() models.KeyValue {
	fmt.Printf("Hello World")
	return models.KeyValue{
		Key:   "This is a random Key",
		Value: "This is a random value",
	}
}
func (p *ProxyService) AddKeyValue() models.KeyValue {
	fmt.Printf("Hello World")
	return models.KeyValue{
		Key:   "This is a random Key",
		Value: "This is a random value",
	}
}
