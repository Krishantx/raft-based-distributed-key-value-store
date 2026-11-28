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
	clientConn            []*grpc.ClientConn
)

func CreateHeartbeatClients() {
	if len(clientConn) == 0 {
		createClientConns()
	}
	HeartbeatClients = make([]pb.HeartbeatServiceClient, len(HOSTNAME))
	for i := 0; i < len(HOSTNAME); i++ {
		HeartbeatClients[i] = pb.NewHeartbeatServiceClient(clientConn[i])
	}
}

func CreateProxyToNodeClients() {
	if len(clientConn) == 0 {
		createClientConns()
	}
	ProxyToNodeClients = make([]pb.ProxyToNodeClient, len(HOSTNAME))
	for i := 0; i < len(HOSTNAME); i++ {
		ProxyToNodeClients[i] = pb.NewProxyToNodeClient(clientConn[i])
	}
}

func CreateRaftOperationsClient() {
	if len(clientConn) == 0 {
		createClientConns()
	}
	RaftOperationsClients = make([]pb.RaftOperationsClient, len(HOSTNAME))
	for i := 0; i < len(HOSTNAME); i++ {
		RaftOperationsClients[i] = pb.NewRaftOperationsClient(clientConn[i])
	}
}

func createClientConns() {
	clientConn = make([]*grpc.ClientConn, len(HOSTNAME))
	for i := 0; i < len(HOSTNAME); i++ {
		clientConn[i] = createClientConn(HOSTNAME[i])
	}
}

func createClientConn(hostname string) *grpc.ClientConn {
	addr := hostname + ":" + PORT
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Failed to connect to client")
	}
	return conn
}
