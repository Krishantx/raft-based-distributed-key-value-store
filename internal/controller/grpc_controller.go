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

type server struct {
	pb.UnimplementedHeartbeatServiceServer
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

func (s *server) ReceiveHeartbeat(context context.Context, in *pb.LeaderDetails) (*pb.ClientConfirmation, error) {
	fmt.Println("Heartbeat Recieved")
	channel <- true
	return &pb.ClientConfirmation{
		Term: 1,
	}, nil
}
func StartgRpcController(config models.Config, ch chan bool) {

	channel = ch

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServ := grpc.NewServer()

	pb.RegisterHeartbeatServiceServer(grpcServ, &server{})
	pb.RegisterVotingServiceServer(grpcServ, &voting{})
	fmt.Println("gRPC Listing on port: 50051")
	err = grpcServ.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
