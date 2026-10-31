package loggers

import (
	"log"
	"os"
)

var (
	Heartbeat *log.Logger
	Raft      *log.Logger
	GRPC      *log.Logger
)

func Init() {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Panicf("Unable to create the directory: %v", err)
	}
	raftFile, err := os.OpenFile("logs/RaftLogs.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 666)
	if err != nil {
		log.Panicf("Cannot Open Log File: %v", err)
	}

	HeartbeatFile, err := os.OpenFile("logs/HeartbeatLogs.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 666)
	if err != nil {
		log.Panicf("Cannot Open Log File: %v", err)
	}

	GRPCFile, err := os.OpenFile("logs/gRpcLogs.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 666)
	if err != nil {
		log.Panicf("Cannot Open Log File: %v", err)
	}

	Heartbeat = log.New(HeartbeatFile, "[Heartbeat]", 0)
	Raft = log.New(raftFile, "[Raft]", 0)
	GRPC = log.New(GRPCFile, "[GRPC]", 0)

}
