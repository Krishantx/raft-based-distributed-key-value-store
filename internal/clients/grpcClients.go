package clients

import (
	"context"
	"log"
	"time"

	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	client pb.VotingServiceClient
	ctx    context.Context
	cancel context.CancelFunc
}

func (this *GRPCClient) GetClinet() pb.VotingServiceClient {
	return this.client
}

func (this *GRPCClient) GetContext() context.Context {
	return this.ctx
}
func (this *GRPCClient) GetCancelFunc() context.CancelFunc {
	return this.cancel
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

func (this *GRPCClientList) GetClienyByIndex(idx int) GRPCClient {
	return this.allClients[idx]
}

func (this *GRPCClientList) ConnectAllNodes() {

	hostname := [5]string{"raft-node-1", "raft-node-2", "raft-node-3", "raft-node-4", "raft-node-5"}

	for i := 0; i < len(this.allClients); i++ {
		this.allClients[i] = NewGRPCClient(hostname[i], "50051")
	}
}

func (this *GRPCClientList) BroadcastAcrossNodes(f func(context.Context, *pb.NodeInfo)) {

}
