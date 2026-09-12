package service

import (
	"context"
	"log"
	models "raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProxyService struct {
	Leader models.Node
}

func NewProxyService() ProxyService {
	node := models.Node{
		Name:     "A",
		Hostname: "raft-node-1",
		Port:     "50051",
	}
	return ProxyService{
		Leader: node,
	}
}

func (p *ProxyService) GetKeyValue(key string) models.KeyValue {
	leaderAddr := p.Leader.Hostname + ":" + p.Leader.Port

	conn, err := grpc.NewClient(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", &leaderAddr, err)
	}

	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.Key{
		Term: 1,
		Key:  key,
	}

	response, err := client.GetKeyValue(ctx, req)
	return models.KeyValue{
		Key:   response.GetKeyValue().Key,
		Value: response.GetKeyValue().Value,
	}
}

func (p *ProxyService) DeleteKeyValue(key string) models.KeyValue {
	leaderAddr := p.Leader.Hostname + ":" + p.Leader.Port

	conn, err := grpc.NewClient(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", &leaderAddr, err)
	}

	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.Key{
		Term: 1,
		Key:  "Key",
	}

	response, err := client.DeleteKeyValue(ctx, req)
	return models.KeyValue{
		Key:   response.GetKeyValue().Key,
		Value: response.GetKeyValue().Value,
	}
}

func (p *ProxyService) PutKeyValue(keyValue models.KeyValue) models.KeyValue {
	leaderAddr := p.Leader.Hostname + ":" + p.Leader.Port
	// Send a gRPC Request to the Leader

	conn, err := grpc.NewClient(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", &leaderAddr, err)
	}

	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.KeyValue{
		Term:  1,
		Key:   keyValue.Key,
		Value: keyValue.Value,
	}

	response, err := client.PutKeyValue(ctx, req)
	return models.KeyValue{
		Key:   response.GetKeyValue().Key,
		Value: response.GetKeyValue().Value,
	}
}

func (p *ProxyService) AddKeyValue(keyValue models.KeyValue) (models.KeyValue, error) {
	leaderAddr := p.Leader.Hostname + ":" + p.Leader.Port
	// Send a gRPC Request to the Leader

	conn, err := grpc.NewClient(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect : %s: %v", leaderAddr, err)
	}

	client := pb.NewProxyToNodeClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.KeyValue{
		Term:  1,
		Key:   keyValue.Key,
		Value: keyValue.Value,
	}

	response, err := client.AddKeyValue(ctx, req)
	if err != nil {
		return models.KeyValue{}, err
	}
	return models.KeyValue{
		Key:   response.GetKeyValue().Key,
		Value: response.GetKeyValue().Value,
	}, nil
}
