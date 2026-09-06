package controller

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	models "raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
)

var channel chan bool

type gRPC_Server struct {
}

func New_gRPC_Server() *gRPC_Server {
	return &gRPC_Server{}
}

type server struct {
	pb.UnimplementedHeartbeatServiceServer
}

type proxyServ struct {
	pb.UnimplementedProxyToNodeServer
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

func (s *proxyServ) GetKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
	fmt.Println("Request recieved for Key : " + req.Key)
	return &pb.Response{
		Result: &pb.Response_KeyValue{
			KeyValue: &pb.KeyValue{
				Term:  1,
				Key:   "abc",
				Value: "cde",
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
	pb.RegisterProxyToNodeServer(grpcServ, &proxyServ{})
	fmt.Println("gRPC Listing on port: 50051")
	err = grpcServ.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
