package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"raft-based-kv/internal/repo"
)

type StartupService struct {
	repository       *repo.Repo
	offsetRepository *[]uint64
}

func NewStartUpService(repository *repo.Repo, offsetRepository *[]uint64) *StartupService {
	return &StartupService{
		repository:       repository,
		offsetRepository: offsetRepository,
	}
}

// Function name is not descriptive enough FIX LATER
func (this *StartupService) ModifyCurrentStateToLatestLog() {
	file, err := os.ReadFile("/app/logs/RaftLogs.log")
	split := bytes.Split(file, []byte("\n"))
	if err != nil {
		log.Fatalf("Cannot Open the log file")
	}
	var offset uint64 = 0
	i := 0
	for _, log := range split {
		(*this.offsetRepository)[i] = offset
		i++
		offset += uint64(len(log)) + 1
		fmt.Println(string(log))
		this.addLogToRepository(log)
	}
}

func (this *StartupService) addLogToRepository(singleLog []byte) {
	logString := string(singleLog)

	decoder := json.NewDecoder(strings.NewReader(logString))
	logEntity := LogEntity{}
	err := decoder.Decode(&logEntity)
	if err != nil {
		log.Printf("Cannot Read String Entity: %s, err: %v", singleLog, err)
	}
	log.Printf("Index: %v, Term: %v, Key: %s", logEntity.Index, logEntity.Term, logEntity.Key)

	switch logEntity.Command {
	case "ADD":
		log.Printf("Added %v", logEntity.Key)
		this.repository.AddKeyValue(logEntity.Key, logEntity.Value)
		break
	case "DELETE":
		log.Printf("DELETEed %v", logEntity.Key)
		this.repository.DeleteKeyValue(logEntity.Key)
		break
	case "PUT":
		log.Printf("Modified %v", logEntity.Key)
		this.repository.UpdateKeyValue(logEntity.Key, logEntity.Value)
		break
	}

}

func (this *StartupService) addToRepository(key string, value string) {
	this.repository.AddKeyValue(key, value)
}

type LogEntity struct {
	Index   uint64
	Term    uint64
	Key     string
	Value   string
	Command string
}
