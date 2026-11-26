package clients

import (
	"log"

	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var HOSTNAME = []string{"raft-node-1", "raft-node-2", "raft-node-3", "raft-node-4", "raft-node-5"}

const PORT = "50051"

var (
	ProxyToNodeClients    []pb.ProxyToNodeClient
	RaftOperationsClients []pb.RaftOperationsClient
	HeartbeatClients      []pb.HeartbeatServiceClient
)

func CreateHeartbeatClients() {
	HeartbeatClients = make([]pb.HeartbeatServiceClient, len(HOSTNAME))
	for i := 0; i < len(HOSTNAME); i++ {
		addr := HOSTNAME[i] + ":" + PORT
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("failed to connect to client")
		}
		HeartbeatClients[i] = pb.NewHeartbeatServiceClient(conn)
	}
}

func CreateProxyToNodeClients() {
	ProxyToNodeClients = make([]pb.ProxyToNodeClient, len(HOSTNAME))
	for i := 0; i < len(HOSTNAME); i++ {
		addr := HOSTNAME[i] + ":" + PORT
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("failed to connect to client")
		}
		ProxyToNodeClients[i] = pb.NewProxyToNodeClient(conn)
	}
}

func CreateRaftOperationsClient() {
	RaftOperationsClients = make([]pb.RaftOperationsClient, len(HOSTNAME))
	for i := 0; i < len(HOSTNAME); i++ {
		addr := HOSTNAME[i] + ":" + PORT
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println("Failed to connect to client")
		}
		RaftOperationsClients[i] = pb.NewRaftOperationsClient(conn)
	}

}
