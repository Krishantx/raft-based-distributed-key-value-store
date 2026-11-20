package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"raft-based-kv/internal/loggers"
	"raft-based-kv/internal/models"
	repository "raft-based-kv/internal/repo"
	pb "raft-based-kv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type leaderRaftOperations struct {
	pb.UnimplementedRaftOperationsServer
}

type proxyToNode struct {
	pb.UnimplementedProxyToNodeServer
	repository   repository.Repo
	raftNodeData *models.NodeRaftData
}

func appendToSelf(logEntry *pb.Entry, nrd *models.NodeRaftData) {
	loggers.Raft.Printf("{\"Index\": \"%v\", \"Term\": \"%v\", \"Command\": \"%s\", \"Key\": \"%s\", \"Value\" : \"%s\"}", logEntry.LastLogIndex+1, logEntry.LastLogTerm, logEntry.Command, logEntry.Key, logEntry.Value)
	nrd.IncrementIndex()
}

func logEntryBuilder(rnd *models.NodeRaftData, c string, k string, v string) pb.Entry {
	return pb.Entry{
		LastLogIndex: rnd.GetIndex(),
		LastLogTerm:  rnd.GetTerm(),
		Command:      c,
		Key:          k,
		Value:        v,
	}
}

func (s *proxyToNode) AddKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
	// Build a pb.Entry obj
	entry := logEntryBuilder(s.raftNodeData, "ADD", req.Key, req.Value)
	appendToSelf(&entry, s.raftNodeData)
	// AppendEntries to Logs of the followers
	err := SendAppendEntryGRPC(&entry, s.raftNodeData)
	// Check if majority Appened
	if err != nil {
		return &pb.Response{}, err
	}

	// CommitEntries to the Logs of followers
	err = SendCommitEntryGRPC()
	if err != nil {
		return &pb.Response{}, err
	}
	log.Println("Commit To Leader Successful")
	s.raftNodeData.IncrementCommitIndex()

	keyValue, err := s.repository.AddKeyValue(req.Key, req.Value)
	if err != nil {
		return &pb.Response{}, err
	}

	return &pb.Response{
		Result: &pb.Response_KeyValue{
			KeyValue: &pb.KeyValue{
				Term:  int32(s.raftNodeData.GetTerm()),
				Key:   keyValue.Key,
				Value: keyValue.Value,
			},
		},
	}, nil
}

func (s *proxyToNode) DeleteKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
	entry := logEntryBuilder(s.raftNodeData, "DELETE", req.Key, "")

	appendToSelf(&entry, s.raftNodeData)

	// AppendEntries to Logs of the followers
	err := SendAppendEntryGRPC(&entry, s.raftNodeData)
	// Check if majority Appened
	if err != nil {
		return &pb.Response{}, err
	}

	// CommitEntries to the Logs of followers
	err = SendCommitEntryGRPC()

	if err != nil {
		return &pb.Response{}, err
	}
	s.raftNodeData.IncrementCommitIndex()
	err = s.repository.DeleteKeyValue(req.Key)
	if err != nil {
		return &pb.Response{}, err
	}
	return &pb.Response{
		Result: &pb.Response_KeyValue{
			KeyValue: &pb.KeyValue{
				Term: int32(s.raftNodeData.GetTerm()),
				Key:  req.Key,
			},
		},
	}, nil
}

func (s *proxyToNode) PutKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
	entry := logEntryBuilder(s.raftNodeData, "PUT", req.Key, req.Value)
	appendToSelf(&entry, s.raftNodeData)

	// AppendEntries to Logs of the followers
	err := SendAppendEntryGRPC(&entry, s.raftNodeData)

	// Check if majority Appened
	if err != nil {
		if err == errors.New("Node is too far behind and needs to catch up") {
			log.Println("======Bring Follower state to fresh data========")

			// Pending Implementation
			UpdateFollowerState()
		}
		return &pb.Response{}, err
	}

	// CommitEntries to the Logs of followers
	err = SendCommitEntryGRPC()

	if err != nil {
		return &pb.Response{}, err
	}
	log.Println("Append To Leader Successful")
	s.raftNodeData.IncrementCommitIndex()

	keyValue, err := s.repository.UpdateKeyValue(req.Key, req.Value)
	fmt.Print("Put Key value")
	if err != nil {
		return &pb.Response{}, err
	}
	return &pb.Response{
		Result: &pb.Response_KeyValue{
			KeyValue: &pb.KeyValue{
				Term:  int32(s.raftNodeData.GetTerm()),
				Key:   keyValue.Key,
				Value: keyValue.Value,
			},
		},
	}, nil
}

func (s *proxyToNode) GetKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
	fmt.Println("get Key value")
	keyValue, err := s.repository.GetKeyValue(req.Key)
	if err != nil {
		return &pb.Response{}, err
	}
	return &pb.Response{
		Result: &pb.Response_KeyValue{
			KeyValue: &pb.KeyValue{
				Term:  int32(s.raftNodeData.GetTerm()),
				Key:   keyValue.Key,
				Value: keyValue.Value,
			},
		},
	}, nil
}
func UpdateFollowerState() {
}

func SendAppendEntryGRPC(entry *pb.Entry, nodeRaftData *models.NodeRaftData) error {
	addrs := "raft-node-"
	cnt := 0
	for i := 1; i < 6; i++ {
		addr := addrs + strconv.Itoa(i) + ":50051"
		if addr != os.Getenv("node_name")+":50051" {

			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Failed to connect : %s: %v", addr, err)
			}

			client := pb.NewRaftOperationsClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			response, err := client.AppendEntries(ctx, entry)

			if err != nil {
				log.Fatalf("Error Sending AppendEntries %s : %v", addr, err)
			} else {
				fmt.Print(response)
				if response.IsAppended {
					cnt++
				} else {
					if response.LastCommitIndex < nodeRaftData.GetCommitIndex() {
						return errors.New("Node is too far behind and needs to catch up")
					}
				}
			}
		}
	}
	if isMajority(cnt + 1) {
		return nil
	}
	return errors.New("Majority nodes did not append")
}

func SendCommitEntryGRPC() error {
	addrs := "raft-node-"
	cnt := 0
	for i := 1; i < 6; i++ {
		addr := addrs + strconv.Itoa(i) + ":50051"
		if addr != os.Getenv("node_name")+":50051" {

			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Failed to connect : %s: %v", addr, err)
			}

			client := pb.NewRaftOperationsClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			req := &pb.Entry{
				LastLogTerm:  1,
				LastLogIndex: 1,
				Command:      "GET",
				Key:          "key",
				Value:        "value",
			}
			response, err := client.CommitEntry(ctx, req)

			if err != nil {
				log.Fatalf("Error Sending AppendEntries %s : %v", addr, err)
			} else {
				fmt.Print(response)
				if response.IsCommited {
					cnt++
				}
			}
		}
	}
	if isMajority(cnt + 1) {
		return nil
	}
	return errors.New("Majority nodes did not Commit")
}
func isMajority(count int) bool {
	totalNodes := 5 //Hardcoded for now, maybe scale later?
	return count > (totalNodes / 2)
}

func StartLeaderController(repository repository.Repo, nodeRaftData *models.NodeRaftData) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("Could not listen on port 50051: %v", err)
	}

	leaderController := grpc.NewServer()
	pb.RegisterProxyToNodeServer(leaderController, &proxyToNode{
		repository:   repository,
		raftNodeData: nodeRaftData,
	})
	pb.RegisterRaftOperationsServer(leaderController, &leaderRaftOperations{})
	fmt.Println("gRPC Listing on port: 50051")
	go func() {
		err = leaderController.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return leaderController, lis
}
