package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"raft-based-kv/internal/enums"
	"raft-based-kv/internal/loggers"
	"raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
	"time"

	grpc "google.golang.org/grpc"
	insecure "google.golang.org/grpc/credentials/insecure"
)

func New_Follower_Service(VotingTimeout int, config *models.Config, timerOps chan enums.TimerOperations) Follower_Service {

	return Follower_Service{
		VotingTimeout: VotingTimeout,
		config:        config,
		timerOps:      timerOps,
	}
}

type Follower_Service struct {
	VotingTimeout int
	config        *models.Config
	timerOps      chan enums.TimerOperations
}

func (this *Follower_Service) StartFollowerService() {
	loggers.Raft.Println("Role Set to Follower")
	go this.StartElectionTimer()

}

func (this *Follower_Service) StartElectionTimer() {
	fmt.Println("Election Timer Started")
	duration := 5 * time.Second
	timer := time.NewTimer(duration)
	for this.config.Role == enums.Follower {
		select {
		//case 1: Leader did not send a Heartbeat, start voting.
		case <-timer.C:
			fmt.Println("Timer expired: Changing role to candidate")
			this.config.Role = enums.Candidate
			break

		case a := <-this.timerOps:
			if a == enums.STOP {
				timer.Stop()
			} else if a == enums.RESET {
				logFile, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 666)
				if err != nil {
					log.Panicln("Could not open log file %s", err)
				}
				defer logFile.Close()
				log.SetOutput(logFile)
				log.Println("Recieved Hear")
				timer.Reset(duration)
			}
		}
	}
}
func grpcClientForVoting(addr string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", addr, err)
	}

	client := pb.NewVotingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.NodeInfo{
		Term: 1,
		Log:  1,
		// NodeName: conf.NodeName,
	}
	response, err := client.StartVoting(ctx, req)

	if err != nil {
		log.Fatalf("Error Sending Heartbear %s : %v", addr, err)
	}
	fmt.Print(response)
}
