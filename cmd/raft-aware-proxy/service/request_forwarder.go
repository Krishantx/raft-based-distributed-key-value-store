package service

import (
	"context"
	"log"
	pb "raft-based-kv/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type nodeDetails struct {
	hostname string
	port     string
}

type clientDetails struct {
	clinet pb.ProxyToNodeClient
	ctx    context.Context
	cancel context.CancelFunc
}

type requestForwarder struct {
	nodeClients [5]clientDetails
}

func NewRequestForwarder() requestForwarder {
	var nodeClients [5]clientDetails
	allNodeDetails := [5]nodeDetails{
		nodeDetails{hostname: "raft-node-1", port: "50051"},
		nodeDetails{hostname: "raft-node-2", port: "50051"},
		nodeDetails{hostname: "raft-node-3", port: "50051"},
		nodeDetails{hostname: "raft-node-4", port: "50051"},
		nodeDetails{hostname: "raft-node-5", port: "50051"},
	}
	for i := 0; i < len(nodeClients); i++ {
		nodeClients[i] = createNewGrpcClient(allNodeDetails[i])
	}
	return requestForwarder{
		nodeClients: nodeClients,
	}
}

func (this *requestForwarder) getKey(key string, req pb.Key) (pb.Response, error) {
	clinet := this.nodeClients[0].clinet
	ctx := this.nodeClients[0].ctx
	response, err := clinet.GetKeyValue(ctx, &req)
	return *response, err
}

func createNewGrpcClient(node nodeDetails) clientDetails {
	addr := node.hostname + ":" + node.port

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}
	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	clientDetails := clientDetails{
		clinet: client,
		ctx:    ctx,
		cancel: cancel,
	}
	return clientDetails
}
