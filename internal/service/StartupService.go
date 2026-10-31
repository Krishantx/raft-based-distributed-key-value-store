package service

import (
	"bytes"
	"log"
	"os"
	"raft-based-kv/internal/repo"
)

type StartupService struct{}

func NewStartUpService() *StartupService {
	return &StartupService{}
}

// Function name is not descriptive enough FIX LATER
func (this *StartupService) ModifyCurrentStateToLatestLog(repository repo.Repo) {
	file, err := os.ReadFile("/app/logs/RaftLogs.log")
	split := bytes.Split(file, []byte("[RAFT]"))
	if err != nil {
		log.Fatalf("Cannot Open the log file")
	}
	for index, element := range split {
		log.Printf("index: %v, element: %v", index, string(element))
	}
}
