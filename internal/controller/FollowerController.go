package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"

	"raft-based-kv/internal/enums"
	"raft-based-kv/internal/loggers"
	"raft-based-kv/internal/models"
	"raft-based-kv/internal/repo"
	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
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

type voting struct {
	pb.UnimplementedVotingServiceServer
	nodeRaftData *models.NodeRaftData
	timerOps     chan enums.TimerOperations
	votingLock   sync.Mutex
}

func (this *voting) StartVoting(ctx context.Context, req *pb.NodeInfo) (*pb.Vote, error) {
	this.votingLock.Lock()
	defer this.votingLock.Unlock()
	fmt.Println("========Starting Voting Procedure=========")
	this.timerOps <- enums.STOP
	log.Println("=========Paused the heartbeat timer to vote========")
	log.Printf("My term: %v, Incoming Term: %v", this.nodeRaftData.GetTerm(), req.Term)
	if req.Term <= int32(this.nodeRaftData.GetTerm()) {
		this.timerOps <- enums.RESET
		return &pb.Vote{
			Vote: false,
		}, nil
	}

	this.nodeRaftData.SetTerm(uint64(req.Term))

	if req.Log < int32(this.nodeRaftData.GetIndex()) {
		this.timerOps <- enums.RESET
		return &pb.Vote{
			Vote: false,
		}, nil
	}
	log.Println("Sending a True vote in return")
	this.timerOps <- enums.RESET
	return &pb.Vote{
		Vote: true,
	}, nil
}

func create_gRPC_clients() [5]pb.RaftOperationsClient {
	var hostnames = [...]string{"raft-node-1", "raft-node-2", "raft-node-3", "raft-node-4", "raft-node-5"}
	var clients [5]pb.RaftOperationsClient
	port := "50051"
	for i := 0; i < len(hostnames); i++ {
		addr := hostnames[i] + ":" + port
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect : %s: %v", addr, err)
		}
		clients[i] = pb.NewRaftOperationsClient(conn)

	}
	return clients
}

func find_leader_client(clients *[5]pb.RaftOperationsClient) (pb.RaftOperationsClient, error) {
	for i := 0; i < len(clients); i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		role, err := clients[i].GetRole(ctx, &emptypb.Empty{})
		if err != nil {
			log.Println("Node returned an error")
		}
		if role.NodeRole == "leader" {
			leader_client := clients[i]
			return leader_client, nil
		}
	}
	return nil, errors.New("Cannot find the leader node")
}

func bulk_append_log_entries() error {
	return nil
}

func followerIndexRecovery(nrd *models.NodeRaftData) {
	clients := create_gRPC_clients()

	// Find who the leader is:
	leader_client, err := find_leader_client(&clients)
	if err != nil {
		return
	}

	log.Printf("Leader Client Found")
	// Ask the leader what its commit index is
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	leader_last_commit_index, err := leader_client.GetLastCommitIndex(ctx, &emptypb.Empty{})
	if err != nil {
		log.Println("Could not fetch the commit count of the leader. Will try again")
		return
	}
	// Start bulk writes of 100

	bulk_append_start_idx := nrd.GetCommitIndex() + 1
	total_logs_to_append := leader_last_commit_index.LeaderLastCommitIndex - bulk_append_start_idx + 1
	var total_iterations float64 = float64(total_logs_to_append / 100)
	total_iterations = math.Ceil(total_iterations)
	for i := uint64(1); i <= uint64(total_iterations); i++ {
		to_commit_index := math.Max(float64(leader_last_commit_index.LeaderLastCommitIndex), float64(bulk_append_start_idx+(i*100)))
		request := &pb.BulkAppendRequest{
			FromCommitIndex: bulk_append_start_idx,
			ToCommitIndex:   uint64(to_commit_index),
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		leader_client.BulkAppendLogEntries(ctx, request)
	}

	err = bulk_append_log_entries()

	if err != nil {
		log.Println("Unable to append log entries")
	}

	err = bulk_commit_log_entries()

	if err != nil {
		log.Println("Unable to commit log entries")
	}
	return
}

func bulk_commit_log_entries() error {
	return nil
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
			IsAppended:      false,
			LastCommitIndex: this.nodeRaftData.GetIndex(),
		}, nil
	}

	if reqLastLogIndex != lastLogIndex {
		log.Println("I have missed few indexes and need to catch up")
		go followerIndexRecovery(this.nodeRaftData)
		return &pb.IsAppended{
			IsAppended:      false,
			LastCommitIndex: this.nodeRaftData.GetIndex(),
		}, nil
	}

	loggers.Raft.Printf("{\"Index\": \"%v\", \"Term\": \"%v\", \"Command\": \"%s\", \"Key\": \"%s\", \"Value\": \"%s\"}", this.nodeRaftData.GetIndex(), this.nodeRaftData.GetTerm(), reqCommand, reqKey, reqValue)

	this.nodeRaftData.IncrementIndex()

	return &pb.IsAppended{
		IsAppended:      true,
		LastCommitIndex: this.nodeRaftData.GetIndex(),
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

	pb.RegisterVotingServiceServer(followerServer, &voting{
		nodeRaftData: nodeRaftData,
		timerOps:     timerOps,
	})
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
