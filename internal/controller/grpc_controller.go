package controller

import (
	"fmt"
	"log"
	"net"
	"raft-based-kv/internal/enums"
	models "raft-based-kv/internal/models"
	repository "raft-based-kv/internal/repo"
	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
)

type gRPC_Server struct {
	repository   repository.Repo
	nodeRaftData *models.NodeRaftData
	timerOps     chan enums.TimerOperations
}

func New_gRPC_Server(repo *repository.Repo, NodeRaftData *models.NodeRaftData, timerOps chan enums.TimerOperations) *gRPC_Server {
	return &gRPC_Server{
		repository:   *repo,
		nodeRaftData: NodeRaftData,
		timerOps:     timerOps,
	}
}

func (s *gRPC_Server) StartgRpcController() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServ := grpc.NewServer()

	pb.RegisterHeartbeatServiceServer(grpcServ, &server{
		timerOps: s.timerOps,
	})
	pb.RegisterVotingServiceServer(grpcServ, &voting{
		nodeRaftData: *s.nodeRaftData,
		timerOps:     s.timerOps,
	})
	pb.RegisterProxyToNodeServer(grpcServ, &proxyServ{
		repository: s.repository,
	})
	fmt.Println("gRPC Listing on port: 50051")
	err = grpcServ.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}

}
