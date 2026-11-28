package models

import (
	"os"
	"strings"

	enums "raft-based-kv/internal/enums"
)

type Config struct {
	NodeName string
	Role     enums.Role
	Follower []string
}

func (this *Config) Init() {
	this.NodeName = os.Getenv("node_name")
	this.Role = enums.Role(os.Getenv("role"))
	this.Follower = strings.Split(os.Getenv("followers"), ",")
}
