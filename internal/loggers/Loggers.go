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
	raftFile, err := os.OpenFile("RaftLogs.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 666)
	if err != nil {
		log.Panicln("Cannot Open Log File")
	}

	HeartbeatFile, err := os.OpenFile("RaftLogs.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 666)
	if err != nil {
		log.Panicln("Cannot Open Log File")
	}

	GRPCFile, err := os.OpenFile("RaftLogs.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 666)
	if err != nil {
		log.Panicln("Cannot Open Log File")
	}

	Heartbeat = log.New(HeartbeatFile, "[Heartbeat]", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	Raft = log.New(raftFile, "[Heartbeat]", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	GRPC = log.New(GRPCFile, "[Heartbeat]", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

}
