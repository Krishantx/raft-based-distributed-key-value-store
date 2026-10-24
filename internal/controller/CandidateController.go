package controller

import (
	"context"
	"fmt"
	"log"
	"net"
	"raft-based-kv/internal/enums"
	models "raft-based-kv/internal/models"
	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
)

type voting struct {
	pb.UnimplementedVotingServiceServer
	nodeRaftData models.NodeRaftData
	timerOps     chan enums.TimerOperations
}

func (this *voting) StartVoting(ctx context.Context, req *pb.NodeInfo) (*pb.Vote, error) {
	fmt.Println("Starting Voting procedure")
	this.timerOps <- enums.RESET
	if req.Term < int32(this.nodeRaftData.GetTerm()) {
		return &pb.Vote{
			Vote: false,
		}, nil
	}

	if req.Log < int32(this.nodeRaftData.GetIndex()) {
		return &pb.Vote{
			Vote: false,
		}, nil
	}

	return &pb.Vote{
		Vote: true,
	}, nil
}

func StartCandidateServer(nodeRaftData models.NodeRaftData, timerOps chan enums.TimerOperations) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("Could Not Listen on port: %v", err)
	}

	controllerServer := grpc.NewServer()

	pb.RegisterVotingServiceServer(controllerServer, &voting{
		nodeRaftData: nodeRaftData,
		timerOps:     timerOps,
	})

	log.Println("controller gRPC Server Listning on port: 50051")

	err = controllerServer.Serve(lis)

	if err != nil {
		log.Panicf("Cannot start gRPC Server: %v", err)
	}
}
