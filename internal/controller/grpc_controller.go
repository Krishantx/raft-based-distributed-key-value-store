package controller

import (
	"context"
	"fmt"
	"log"
	"net"
	models "raft-based-kv/internal/models"
	repository "raft-based-kv/internal/repo"
	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
)

var channel chan bool

type gRPC_Server struct {
	repository repository.Repo
}

func New_gRPC_Server(repo *repository.Repo) *gRPC_Server {
	return &gRPC_Server{
		repository: *repo,
	}
}

type server struct {
	pb.UnimplementedHeartbeatServiceServer
}

type proxyServ struct {
	pb.UnimplementedProxyToNodeServer
	repository repository.Repo
}

type voting struct {
	pb.UnimplementedVotingServiceServer
}

func (s *voting) StartVoting(ctx context.Context, req *pb.NodeInfo) (*pb.Vote, error) {
	fmt.Println("Starting Voting procedure")

	return &pb.Vote{
		Vote: false,
	}, nil
}

func (s *proxyServ) DeleteKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
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

func (s *proxyServ) PutKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
	keyValue, err := s.repository.UpdateKeyValue(req.Key, req.Value)
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

func (s *proxyServ) GetKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
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
func (s *proxyServ) AddKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
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
func (s *server) ReceiveHeartbeat(context context.Context, in *pb.Leader) (*pb.ClientConfirmation, error) {
	// fmt.Println("Heartbeat Recieved")
	channel <- true
	return &pb.ClientConfirmation{
		Term: 1,
	}, nil
}
func (s *gRPC_Server) StartgRpcController(config models.Config, ch chan bool, voteChan chan bool) {

	channel = ch

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServ := grpc.NewServer()

	pb.RegisterHeartbeatServiceServer(grpcServ, &server{})
	pb.RegisterVotingServiceServer(grpcServ, &voting{})
	pb.RegisterProxyToNodeServer(grpcServ, &proxyServ{
		repository: s.repository,
	})
	fmt.Println("gRPC Listing on port: 50051")
	err = grpcServ.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
