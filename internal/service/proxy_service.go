package service

import (
	"fmt"
	models "raft-based-kv/internal/models"
)

type ProxyService struct {
	Leader models.Node
}

func (p ProxyService) addKeyValue() any {
	panic("unimplemented")
}

func NewProxyService() ProxyService {
	node := models.Node{
		Name:     "A",
		Hostname: "raft-node-a",
	}
	return ProxyService{
		Leader: node,
	}
}

func (p *ProxyService) GetKeyValue() models.KeyValue {
	fmt.Printf("Hello World")
	return models.KeyValue{
		Key:   "This is a random Key",
		Value: "This is a random value",
	}
}

func (p *ProxyService) PutKeyValue() models.KeyValue {
	fmt.Printf("Hello World")
	return models.KeyValue{
		Key:   "This is a random Key",
		Value: "This is a random value",
	}
}
func (p *ProxyService) DeleteKeyValue() models.KeyValue {
	fmt.Printf("Hello World")
	return models.KeyValue{
		Key:   "This is a random Key",
		Value: "This is a random value",
	}
}
func (p *ProxyService) AddKeyValue() models.KeyValue {
	fmt.Printf("Hello World")
	return models.KeyValue{
		Key:   "This is a random Key",
		Value: "This is a random value",
	}
}
