package helpers

import (
	"context"
	"errors"
	"log"
	"time"

	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetLeaderClient() (pb.RaftOperationsClient, error) {
	clients := CreateRaftOperationsClients()
	for i := 0; i < len(clients); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		role, err := clients[i].GetRole(ctx, &emptypb.Empty{})
		if err != nil {
			log.Println("Node returned an error")
		}
		if role.NodeRole == "leader" {
			leader_client := clients[i]
			return leader_client, nil
		}
	}
	return nil, errors.New("Cannot find the leader node")
}

func CreateRaftOperationsClients() [5]pb.RaftOperationsClient {
	var hostnames = [...]string{"raft-node-1", "raft-node-2", "raft-node-3", "raft-node-4", "raft-node-5"}
	var clients [5]pb.RaftOperationsClient
	port := "50051"
	for i := 0; i < len(hostnames); i++ {
		addr := hostnames[i] + ":" + port
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect : %s: %v", addr, err)
		}
		clients[i] = pb.NewRaftOperationsClient(conn)

	}
	return clients
}
