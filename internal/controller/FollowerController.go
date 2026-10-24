package controller

import (
	"context"
	"fmt"
	"log"
	"net"
	"raft-based-kv/internal/enums"
	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHeartbeatServiceServer
	timerOps chan enums.TimerOperations
}

func (s *server) ReceiveHeartbeat(context context.Context, in *pb.Leader) (*pb.ClientConfirmation, error) {
	fmt.Println("Heartbeat Recieved")
	s.timerOps <- enums.RESET
	return &pb.ClientConfirmation{
		Term: 1,
	}, nil
}

func StartFollowerController(timerOps chan enums.TimerOperations) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("Cannot listen on port 50051: %v", err)
	}

	followerServer := grpc.NewServer()

	pb.RegisterHeartbeatServiceServer(followerServer, &server{
		timerOps: timerOps,
	})
	go func() {
		err = followerServer.Serve(lis)
		if err != nil {
			log.Fatalf("Error creating a gRPC server: %v", err)
		}
	}()

	return followerServer, lis
}
