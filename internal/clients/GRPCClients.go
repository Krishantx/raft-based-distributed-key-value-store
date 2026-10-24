package clients

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	client any //Fix this later
	ctx    context.Context
	cancel context.CancelFunc
}

func (this *GRPCClient) GetClinet() any {
	return this.client
}

func (this *GRPCClient) GetContext() context.Context {
	return this.ctx
}
func (this *GRPCClient) GetCancelFunc() context.CancelFunc {
	return this.cancel
}
func NewGRPCClient[T any](hostname string, port string, newClient func(grpc.ClientConnInterface) T) GRPCClient {
	addr := hostname + ":" + port
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := newClient(conn)
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
