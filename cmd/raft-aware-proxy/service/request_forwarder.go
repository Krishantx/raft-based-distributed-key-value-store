package service

import (
	"context"
	"fmt"
	"log"
	pb "raft-based-kv/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	HOSTNAME string = "raft-node-1"
	PORT     string = "50051"
)

type nodeDetails struct {
	hostname string
	port     string
}

type requestForwarder struct {
	allNodeDetails [5]nodeDetails
	clients        [5]pb.ProxyToNodeClient
}

func NewRequestForwarder() requestForwarder {
	allNodeDetails := [5]nodeDetails{
		nodeDetails{hostname: "raft-node-1", port: "50051"},
		nodeDetails{hostname: "raft-node-2", port: "50051"},
		nodeDetails{hostname: "raft-node-3", port: "50051"},
		nodeDetails{hostname: "raft-node-4", port: "50051"},
		nodeDetails{hostname: "raft-node-5", port: "50051"},
	}

	return requestForwarder{
		allNodeDetails: allNodeDetails,
	}
}

func (this *requestForwarder) getKeyValue(key string, req pb.Key) (pb.Response, error) {
	addr := HOSTNAME + ":" + PORT
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Println("Request Reached the forwarder")
	response, err := client.GetKeyValue(ctx, &req)
	fmt.Printf("Request reached the request forwarder")
	return *response, err
}

func (this *requestForwarder) deleteKeyValue(key string, req pb.Key) (pb.Response, error) {
	addr := HOSTNAME + ":" + PORT
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Println("Request Reached the forwarder")
	response, err := client.DeleteKeyValue(ctx, &req)
	return *response, err
}

func (this *requestForwarder) addKeyValue(req pb.KeyValue) (pb.Response, error) {
	addr := HOSTNAME + ":" + PORT
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Println("Request Reached the forwarder")
	response, err := client.AddKeyValue(ctx, &req)
	return *response, err
}

func (this *requestForwarder) putKeyValue(req pb.KeyValue) (pb.Response, error) {
	addr := HOSTNAME + ":" + PORT
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Println("Request Reached the forwarder")
	response, err := client.PutKeyValue(ctx, &req)
	return *response, err
}

func (this *requestForwarder) broadcastAllNodes() {
	var clients [5]pb.ProxyToNodeClient
	for i := 0; i < 5; i++ {
		addr := HOSTNAME + ":" + PORT
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect : %s: %v", addr, err)
		}

		client := pb.NewProxyToNodeClient(conn)
		clients[i] = client
	}
}
