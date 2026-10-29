package models

import (
	"fmt"
	"os"
	enums "raft-based-kv/internal/enums"
	"strconv"
	"strings"
)

type Config struct {
	NodeName     string
	ElectionTime int
	Role         enums.Role
	Follower     []string
}

func (this *Config) Init() {

	this.NodeName = os.Getenv("node_name")
	t := os.Getenv("election_time")
	fmt.Println("ElectionTime :" + t)
	time, err := strconv.Atoi(os.Getenv("election_time"))
	if err != nil {
		fmt.Printf("Cannot parse Election Time: %s", err)
	}
	this.ElectionTime = time
	this.Role = enums.Role(os.Getenv("role"))
	this.Follower = strings.Split(os.Getenv("followers"), ",")
}
