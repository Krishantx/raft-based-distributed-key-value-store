package controller

import (
	"context"
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

type server struct {
	pb.UnimplementedHeartbeatServiceServer
	timerOps chan enums.TimerOperations
}

type proxyServ struct {
	pb.UnimplementedProxyToNodeServer
	repository repository.Repo
}

type voting struct {
	pb.UnimplementedVotingServiceServer
	nodeRaftData models.NodeRaftData
	timerOps     chan enums.TimerOperations
}

func (this *voting) StartVoting(ctx context.Context, req *pb.NodeInfo) (*pb.Vote, error) {
	fmt.Println("Starting Voting procedure")
	// Have a channel that lets the follower service know to stop the heartbeat timer.
	// That way the follower service have maintain the ownership of the heartbeat timer.
	this.timerOps <- enums.STOP
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

func (s *proxyServ) DeleteKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
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

func (s *proxyServ) PutKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
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

func (s *proxyServ) GetKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
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

func (s *proxyServ) AddKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
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
func (s *server) ReceiveHeartbeat(context context.Context, in *pb.Leader) (*pb.ClientConfirmation, error) {
	fmt.Println("Heartbeat Recieved")
	s.timerOps <- enums.RESET
	return &pb.ClientConfirmation{
		Term: 1,
	}, nil
}
func (s *gRPC_Server) StartgRpcController(config models.Config) {
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
