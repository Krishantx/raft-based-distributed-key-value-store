package service

import (
	"context"
	"errors"
	"net"
	"sync"

	"log"
	"time"

	"raft-based-kv/internal/clients"
	"raft-based-kv/internal/enums"
	"raft-based-kv/internal/helpers"
	"raft-based-kv/internal/models"
	pb "raft-based-kv/proto"

	repository "raft-based-kv/internal/repo"

	"google.golang.org/grpc"
)

// Send Heartbeat every 100ms

type LeaderService struct {
	node_name    string
	term         int
	conf         *models.Config
	NodeRaftData *models.NodeRaftData
	repo         repository.Repo
	controller   *grpc.Server
	lis          net.Listener
}

func NewLeaderService(config *models.Config, nodeRaftData *models.NodeRaftData, repo repository.Repo) LeaderService {
	return LeaderService{
		node_name:    config.NodeName,
		term:         1,
		conf:         config,
		NodeRaftData: nodeRaftData,
		repo:         repo,
	}
}

type Counter struct {
	count int
	mu    sync.Mutex
	wg    sync.WaitGroup
}

func (c *Counter) incrementCount() {
	c.mu.Lock()
	c.count = c.count + 1
	c.mu.Unlock()
}

func (l *LeaderService) appendToAllFollowers(entry *pb.Entry) int {
	clients := clients.RaftOperationsClients
	count := Counter{
		count: 1,
	}
	for i := 0; i < len(clients); i++ {
		count.wg.Go(func() {
			l.appendToFollower(clients[i], &count, entry)
		})
	}
	count.wg.Wait()

	return count.count
}

func (l *LeaderService) appendToFollower(client pb.RaftOperationsClient, counter *Counter, entry *pb.Entry) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	isAppended, err := client.AppendEntries(ctx, entry)

	if err != nil {
		log.Printf("Cannot Append Entries to the follower: %v", err)
		return
	}

	if isAppended == nil {
		log.Printf("isAppended is nil and cannot be confirmed")
		return
	}
	if !isAppended.IsAppended {
		return
	}

	if isAppended.IsAppended {
		counter.incrementCount()
	}
}

func (l *LeaderService) commitToAllFollowers(entry *pb.Entry) int {
	clients := clients.RaftOperationsClients
	count := Counter{
		count: 0,
	}
	for i := 0; i < len(clients); i++ {
		count.wg.Go(func() {
			l.commitToFollower(&count, clients[i], entry)
		})
	}
	count.wg.Wait()
	return count.count
}

func (l *LeaderService) commitToFollower(count *Counter, client pb.RaftOperationsClient, entry *pb.Entry) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	isCommited, err := client.CommitEntry(ctx, entry)

	if err != nil {
		log.Printf("Error Commiting Log Entry: %v", err)
		return
	}
	if isCommited == nil {
		log.Printf("Error is Commited payload is nil")
		return
	}

	if !isCommited.IsCommited {
		return
	}
	count.incrementCount()
}
func (l *LeaderService) GetKeyValue(key string) (models.KeyValue, error) {
	return l.repo.GetKeyValue(key)
}

func (l *LeaderService) DeleteKeyValue(key string) error {
	return l.repo.DeleteKeyValue(key)
}

func (l *LeaderService) AddKeyValue(key string, value string) (models.KeyValue, error) {
	entry := helpers.LogEntryBuilder(l.NodeRaftData, "ADD", key, value)
	helpers.AppendToSelf(&entry, l.NodeRaftData)
	cnt := l.appendToAllFollowers(&entry)

	if cnt <= 3 {
		return models.KeyValue{}, errors.New("Could Not be appended by majority nodes")
	}

	cnt = l.commitToAllFollowers(&entry)
	if cnt >= 2 {
		// Commit to self
		l.repo.AddKeyValue(key, value)
		l.NodeRaftData.IncrementCommitIndex()
	}

	return l.repo.GetKeyValue(key)
}

func (l *LeaderService) UpdateKeyValue(key string, value string) (models.KeyValue, error) {
	return l.repo.UpdateKeyValue(key, value)
}

func (this *LeaderService) grpcClient(client pb.HeartbeatServiceClient) {

	req := &pb.Leader{NodeName: this.node_name}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := client.ReceiveHeartbeat(ctx, req)
	if err != nil {
		log.Printf("Error Sending Heartbeart: %v", err)
	} else {
		if int(response.Term) == this.term {
			// fmt.Println("The current Term is correct and I am still the leader")
		} else {
			// fmt.Println("The term does not match and I am no longer the leader")
		}
	}
}

func (this *LeaderService) SendHeartbeat() {
	clients.CreateHeartbeatClients()
	for i := 0; i < 5; i++ {
		go this.grpcClient(clients.HeartbeatClients[i])
	}
}

func (l *LeaderService) HeartbeatLoop() {
	log.Println("Sendign recurrring heartbeats")
	for l.conf.Role == enums.Leader {
		time.Sleep(time.Millisecond * 500)
		l.SendHeartbeat()
	}
}
