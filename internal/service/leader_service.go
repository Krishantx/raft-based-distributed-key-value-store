package service

import (
	"context"
	"fmt"
	"strconv"

	controller "raft-based-kv/internal/controller"
	"raft-based-kv/internal/enums"
	// "fmt"
	"log"
	"raft-based-kv/internal/loggers"
	"raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
	"time"

	repository "raft-based-kv/internal/repo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Send Heartbeat every 100ms

type LeaderService struct {
	node_name    string
	term         int
	conf         models.Config
	timerOps     chan enums.TimerOperations
	nodeRaftData models.NodeRaftData
}

func NewLeaderService(config models.Config, TimerOps chan enums.TimerOperations, nodeRaftData models.NodeRaftData) LeaderService {
	return LeaderService{
		node_name:    config.NodeName,
		term:         1,
		conf:         config,
		timerOps:     TimerOps,
		nodeRaftData: nodeRaftData,
	}
}

func (s *LeaderService) StartLeaderService() {

	Repo := repository.NewRepository()
	grpc_Server := controller.New_gRPC_Server(Repo, &s.nodeRaftData, s.timerOps)
	go grpc_Server.StartgRpcController()

	loggers.Raft.Println("Role Set to Leader")
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		s.SendHeartbeat()
	}
}

func (this *LeaderService) grpcClient(addr string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewHeartbeatServiceClient(conn)

	req := &pb.Leader{NodeName: this.node_name}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := client.ReceiveHeartbeat(ctx, req)

	if err != nil {
		log.Println("Error Sending Heartbear %s : %v", addr, err)
	}

	if int(response.Term) == this.term {
		// fmt.Println("The current Term is correct and I am still the leader")
	} else {
		// fmt.Println("The term does not match and I am no longer the leader")
	}
}

func (this *LeaderService) SendHeartbeat() {
	addrs := "raft-node-"
	for i := 1; i < 6; i++ {
		addr := addrs + strconv.Itoa(i) + ":50051"
		fmt.Println(addr + " : " + this.node_name + ":50051")
		if addr != this.node_name+":50051" {
			go this.grpcClient(addr)
		}
	}
}
