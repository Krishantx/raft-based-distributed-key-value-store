package controller

import (
	"context"
	"fmt"
	"log"
	"net"
	"raft-based-kv/internal/enums"
	"raft-based-kv/internal/loggers"
	"raft-based-kv/internal/models"
	"raft-based-kv/internal/repo"
	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHeartbeatServiceServer
	timerOps chan enums.TimerOperations
}

type raftOperations struct {
	pb.UnimplementedRaftOperationsServer
	nodeRaftData *models.NodeRaftData
	repository   *repo.Repo
}

func (this *raftOperations) AppendEntries(contex context.Context, req *pb.Entry) (*pb.IsAppended, error) {
	lastLogIndex := this.nodeRaftData.GetIndex()
	lastLogTerm := this.nodeRaftData.GetTerm()
	// CAN BE GET, ADD, MODIFY, DELETE
	reqCommand := req.Command
	reqKey := req.Key
	reqValue := req.Value
	reqLastLogIndex := req.LastLogIndex
	reqLastLogTerm := req.LastLogTerm

	if reqLastLogTerm < lastLogTerm {
		log.Println("The Request Comes from a stale leader and is ignored")
		return &pb.IsAppended{
			IsAppended: false,
		}, nil
	}

	if reqLastLogIndex != lastLogIndex {
		log.Println("I have missed few indexes and need to catch up")
		return &pb.IsAppended{
			IsAppended: false,
		}, nil
	}

	loggers.Raft.Printf("Index: %v, Term: %v, Command: %s, Key: %s, Value: %s", this.nodeRaftData.GetIndex(), this.nodeRaftData.GetTerm(), reqCommand, reqKey, reqValue)

	this.nodeRaftData.IncrementIndex()

	return &pb.IsAppended{
		IsAppended: true,
	}, nil
}

func (this *raftOperations) CommitEntry(context context.Context, req *pb.Entry) (*pb.IsCommited, error) {

	this.nodeRaftData.IncrementCommitIndex()
	return &pb.IsCommited{
		IsCommited: true,
	}, nil
}

func (s *server) ReceiveHeartbeat(context context.Context, in *pb.Leader) (*pb.ClientConfirmation, error) {
	fmt.Println("Heartbeat Recieved")
	s.timerOps <- enums.RESET
	return &pb.ClientConfirmation{
		Term: 1,
	}, nil
}

func StartFollowerController(timerOps chan enums.TimerOperations, nodeRaftData *models.NodeRaftData) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("Cannot listen on port 50051: %v", err)
	}

	followerServer := grpc.NewServer()

	pb.RegisterHeartbeatServiceServer(followerServer, &server{
		timerOps: timerOps,
	})
	pb.RegisterRaftOperationsServer(followerServer, &raftOperations{
		nodeRaftData: nodeRaftData,
	})
	go func() {
		err = followerServer.Serve(lis)
		if err != nil {
			log.Fatalf("Error creating a gRPC server: %v", err)
		}
	}()

	return followerServer, lis
}
