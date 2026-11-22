package helpers

import (
	"raft-based-kv/internal/loggers"
	"raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
)

func AppendToSelf(logEntry *pb.Entry, nrd *models.NodeRaftData) {
	loggers.Raft.Printf("{\"Index\": \"%v\", \"Term\": \"%v\", \"Command\": \"%s\", \"Key\": \"%s\", \"Value\" : \"%s\"}", logEntry.LogIndex, logEntry.LogTerm, logEntry.Command, logEntry.Key, logEntry.Value)
	nrd.IncrementIndex()
}

func LogEntryBuilder(rnd *models.NodeRaftData, c string, k string, v string) pb.Entry {
	return pb.Entry{
		LogIndex: rnd.GetIndex(),
		LogTerm:  rnd.GetTerm(),
		Command:  c,
		Key:      k,
		Value:    v,
	}
}
