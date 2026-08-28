package controller

import (
	"fmt"
	"log"
	"net"

	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSendHeartbeatServer
}

func StartgRpcController() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServ := grpc.NewServer()

	pb.RegisterSendHeartbeatServer(grpcServ, &server{})

	fmt.Println("gRPC Listing on port: 50051")
	err = grpcServ.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
