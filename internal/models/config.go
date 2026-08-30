package models

// Change Role type to enum later
type Config struct {
	NodeName     string
	ElectionTime int
	Role         string
	Follower     []string
}
