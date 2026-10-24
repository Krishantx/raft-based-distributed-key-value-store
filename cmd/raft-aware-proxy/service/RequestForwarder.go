package service

import (
	"context"
	"errors"
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
	clients := this.broadcastAllNodes()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 0; i < 5; i++ {
		response, err := clients[i].GetKeyValue(ctx, &req)
		if err != nil {
			fmt.Println(err)
		}
		if response != nil {
			return *response, err
		}
	}
	err := errors.New("None of the Nodes returned a response")
	return pb.Response{}, err
}

func (this *requestForwarder) deleteKeyValue(key string, req pb.Key) (pb.Response, error) {
	clients := this.broadcastAllNodes()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 0; i < 5; i++ {
		response, err := clients[i].DeleteKeyValue(ctx, &req)
		if err != nil {
			fmt.Println(err)
		}
		if response != nil {
			return *response, err
		}
	}
	err := errors.New("None of the Nodes returned a response")
	return pb.Response{}, err
}

func (this *requestForwarder) addKeyValue(req pb.KeyValue) (pb.Response, error) {
	clients := this.broadcastAllNodes()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 0; i < 5; i++ {
		response, err := clients[i].AddKeyValue(ctx, &req)
		if err != nil {
			fmt.Println(err)
		}
		if response != nil {
			return *response, err
		}
	}
	err := errors.New("None of the Nodes returned a response")
	return pb.Response{}, err
}

func (this *requestForwarder) putKeyValue(req pb.KeyValue) (pb.Response, error) {
	clients := this.broadcastAllNodes()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 0; i < 5; i++ {
		response, err := clients[i].PutKeyValue(ctx, &req)
		if err != nil {
			fmt.Println(err)
		}
		if response != nil {
			return *response, err
		}
	}
	err := errors.New("None of the Nodes returned a response")
	return pb.Response{}, err
}

func (this *requestForwarder) broadcastAllNodes() [5]pb.ProxyToNodeClient {
	var clients [5]pb.ProxyToNodeClient
	for i := 0; i < 5; i++ {
		addr := this.allNodeDetails[i].hostname + ":" + this.allNodeDetails[i].port
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect : %s: %v", addr, err)
		}

		client := pb.NewProxyToNodeClient(conn)
		clients[i] = client
	}
	return clients
}
