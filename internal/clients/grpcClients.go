package clients

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "raft-based-kv/proto"
)

type GRPCClient struct {
	client pb.VotingServiceClient
	ctx    context.Context
	cancel context.CancelFunc
}

func NewGRPCClient(hostname string, port string) GRPCClient {
	addr := hostname + ":" + port
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewVotingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	grpcClient := GRPCClient{
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}
	return grpcClient
}

type GRPCClientList struct {
	allClients [5]GRPCClient
}
