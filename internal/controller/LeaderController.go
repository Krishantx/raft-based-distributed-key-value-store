package controller

import (
	"context"
	"fmt"
	"log"
	"net"
	repository "raft-based-kv/internal/repo"
	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
)

type proxyToNode struct {
	pb.UnimplementedProxyToNodeServer
	repository repository.Repo
}

func (s *proxyToNode) DeleteKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
	fmt.Print("Delete Key value")
	err := s.repository.DeleteKeyValue(req.Key)
	if err != nil {
		return &pb.Response{}, err
	}
	return &pb.Response{
		Result: &pb.Response_KeyValue{
			KeyValue: &pb.KeyValue{
				Term: 1,
				Key:  req.Key,
			},
		},
	}, nil
}

func (s *proxyToNode) PutKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
	keyValue, err := s.repository.UpdateKeyValue(req.Key, req.Value)
	fmt.Print("Put Key value")
	if err != nil {
		return &pb.Response{}, err
	}
	return &pb.Response{
		Result: &pb.Response_KeyValue{
			KeyValue: &pb.KeyValue{
				Term:  1,
				Key:   keyValue.Key,
				Value: keyValue.Value,
			},
		},
	}, nil
}

func (s *proxyToNode) GetKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
	fmt.Println("get Key value")
	keyValue, err := s.repository.GetKeyValue(req.Key)
	if err != nil {
		return &pb.Response{}, err
	}
	return &pb.Response{
		Result: &pb.Response_KeyValue{
			KeyValue: &pb.KeyValue{
				Term:  1,
				Key:   keyValue.Key,
				Value: keyValue.Value,
			},
		},
	}, nil
}

func (s *proxyToNode) AddKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
	fmt.Print("add Key value")
	keyValue, err := s.repository.AddKeyValue(req.Key, req.Value)
	if err != nil {

		return &pb.Response{}, err
	}
	return &pb.Response{
		Result: &pb.Response_KeyValue{
			KeyValue: &pb.KeyValue{
				Term:  1,
				Key:   keyValue.Key,
				Value: keyValue.Value,
			},
		},
	}, nil
}

func StartLeaderController(repository repository.Repo) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("Could not listen on port 50051: %v", err)
	}

	leaderController := grpc.NewServer()

	pb.RegisterProxyToNodeServer(leaderController, &proxyToNode{
		repository: repository,
	})
	fmt.Println("gRPC Listing on port: 50051")
	err = leaderController.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}

}
