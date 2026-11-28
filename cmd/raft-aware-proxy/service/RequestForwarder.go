package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	clnt "raft-based-kv/internal/clients"
	pb "raft-based-kv/proto"
)

type requestForwarder struct{}

func NewRequestForwarder() *requestForwarder {
	clnt.CreateProxyToNodeClients()
	return &requestForwarder{}
}

func (this *requestForwarder) getKeyValue(key string, req *pb.Key) (*pb.Response, error) {
	clients := clnt.ProxyToNodeClients
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := clients[i].GetKeyValue(ctx, req)
		if err != nil {
			fmt.Println(err)
		}
		if response != nil {
			return response, err
		}
	}
	err := errors.New("None of the Nodes returned a response")
	return &pb.Response{}, err
}

func (this *requestForwarder) deleteKeyValue(key string, req *pb.Key) (*pb.Response, error) {
	clients := clnt.ProxyToNodeClients
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := clients[i].DeleteKeyValue(ctx, req)
		if err != nil {
			fmt.Println(err)
		}
		if response != nil {
			return response, err
		}
	}
	err := errors.New("None of the Nodes returned a response")
	return &pb.Response{}, err
}

func (this *requestForwarder) addKeyValue(req *pb.KeyValue) (*pb.Response, error) {
	clients := clnt.ProxyToNodeClients
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		response, err := clients[i].AddKeyValue(ctx, req)
		if err != nil {
			fmt.Println(err)
		}
		if response != nil {
			return response, err
		}
	}
	err := errors.New("None of the Nodes returned a response")
	return &pb.Response{}, err
}

func (this *requestForwarder) putKeyValue(req *pb.KeyValue) (*pb.Response, error) {
	clients := clnt.ProxyToNodeClients
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		response, err := clients[i].PutKeyValue(ctx, req)
		if err != nil {
			fmt.Println(err)
		}
		if response != nil {
			return response, err
		}
	}
	err := errors.New("None of the Nodes returned a response")
	return &pb.Response{}, err
}
