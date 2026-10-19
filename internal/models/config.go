package models

import (
	enums "raft-based-kv/internal/enums"
)

// Change Role type to enum later
type Config struct {
	NodeName     string
	ElectionTime int
	Role         enums.Role
	Follower     []string
}
