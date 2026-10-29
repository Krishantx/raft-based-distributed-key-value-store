package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"raft-based-kv/internal/loggers"
	"raft-based-kv/internal/models"
	repository "raft-based-kv/internal/repo"
	pb "raft-based-kv/proto"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type proxyToNode struct {
	pb.UnimplementedProxyToNodeServer
	repository   repository.Repo
	raftNodeData models.NodeRaftData
}

func (s *proxyToNode) DeleteKeyValue(ctx context.Context, req *pb.Key) (*pb.Response, error) {
	index, term := s.raftNodeData.GetIndex(), s.raftNodeData.GetTerm()
	command := "DELETE"
	key := req.Key
	entry := pb.Entry{
		LastLogIndex: s.raftNodeData.GetIndex(),
		LastLogTerm:  s.raftNodeData.GetTerm(),
		Command:      command,
		Key:          req.Key,
	}
	loggers.Raft.Printf("Index: %v, Term: %v, Command: %s, Key: %s", index, term, command, key)
	s.raftNodeData.IncrementIndex()
	log.Printf("Index: %v, Term: %v, Command: %s, Key: %s", index, term, command, key)
	// AppendEntries to Logs of the followers
	err := SendAppendEntryGRPC(entry)

	// Check if majority Appened
	if err != nil {
		return &pb.Response{}, err
	}

	// CommitEntries to the Logs of followers
	err = SendCommitEntryGRPC()

	if err != nil {
		return &pb.Response{}, err
	}
	log.Println("Append To Leader Successful")
	s.raftNodeData.IncrementCommitIndex()
	fmt.Print("Delete Key value")
	err = s.repository.DeleteKeyValue(req.Key)
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

func (s *proxyToNode) PutKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
	index, term := s.raftNodeData.GetIndex(), s.raftNodeData.GetTerm()
	command := "PUT"
	key := req.Key
	value := req.Value
	entry := pb.Entry{
		LastLogIndex: s.raftNodeData.GetIndex(),
		LastLogTerm:  s.raftNodeData.GetTerm(),
		Command:      command,
		Key:          req.Key,
		Value:        req.Value,
	}
	loggers.Raft.Printf("Index: %v, Term: %v, Command: %s, Key: %s, Value: %s", index, term, command, key, value)
	s.raftNodeData.IncrementIndex()
	log.Printf("Index: %v, Term: %v, Command: %s, Key: %s, Value: %s", index, term, command, key, value)
	// AppendEntries to Logs of the followers
	err := SendAppendEntryGRPC(entry)

	// Check if majority Appened
	if err != nil {
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
				Term:  1,
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
				Term:  1,
				Key:   keyValue.Key,
				Value: keyValue.Value,
			},
		},
	}, nil
}
func SendAppendEntryGRPC(entry pb.Entry) error {
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

			response, err := client.AppendEntries(ctx, &entry)

			if err != nil {
				log.Fatalf("Error Sending AppendEntries %s : %v", addr, err)
			} else {
				fmt.Print(response)
				if response.IsAppended {
					cnt++
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

func (s *proxyToNode) AddKeyValue(ctx context.Context, req *pb.KeyValue) (*pb.Response, error) {
	//Step 1 Append Entries to your own log
	index, term := s.raftNodeData.GetIndex(), s.raftNodeData.GetTerm()
	command := "ADD"
	key := req.Key
	value := req.Value
	entry := pb.Entry{
		LastLogIndex: s.raftNodeData.GetIndex(),
		LastLogTerm:  s.raftNodeData.GetTerm(),
		Command:      command,
		Key:          req.Key,
		Value:        req.Value,
	}
	loggers.Raft.Printf("Index: %v, Term: %v, Command: %s, Key: %s, Value: %s", index, term, command, key, value)
	s.raftNodeData.IncrementIndex()
	log.Printf("Index: %v, Term: %v, Command: %s, Key: %s, Value: %s", index, term, command, key, value)
	// AppendEntries to Logs of the followers
	err := SendAppendEntryGRPC(entry)

	// Check if majority Appened
	if err != nil {
		return &pb.Response{}, err
	}

	// CommitEntries to the Logs of followers
	err = SendCommitEntryGRPC()

	if err != nil {
		return &pb.Response{}, err
	}
	log.Println("Append To Leader Successful")
	s.raftNodeData.IncrementCommitIndex()

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

func StartLeaderController(repository repository.Repo, nodeRaftData models.NodeRaftData) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("Could not listen on port 50051: %v", err)
	}

	leaderController := grpc.NewServer()
	pb.RegisterProxyToNodeServer(leaderController, &proxyToNode{
		repository:   repository,
		raftNodeData: nodeRaftData,
	})

	fmt.Println("gRPC Listing on port: 50051")
	go func() {
		err = leaderController.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return leaderController, lis
}
