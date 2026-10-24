package service

import (
	models "raft-based-kv/internal/models"
	pb "raft-based-kv/proto"
)

type ProxyService struct {
	requestForwarder requestForwarder
}

func NewProxyService() ProxyService {
	requestForwarder := NewRequestForwarder()
	return ProxyService{
		requestForwarder: requestForwarder,
	}
}

func (p *ProxyService) GetKeyValue(key string) (models.KeyValue, error) {
	req := &pb.Key{
		Term: 1,
		Key:  key,
	}

	response, err := p.requestForwarder.getKeyValue(key, *req)

	if err != nil {
		return models.KeyValue{}, err
	}

	return models.KeyValue{
		Key:   response.GetKeyValue().Key,
		Value: response.GetKeyValue().Value,
	}, nil
}

func (p *ProxyService) DeleteKeyValue(key string) (models.KeyValue, error) {
	req := &pb.Key{
		Term: 1,
		Key:  key,
	}

	response, err := p.requestForwarder.deleteKeyValue(key, *req)

	if err != nil {
		return models.KeyValue{}, err
	}

	return models.KeyValue{
		Key:   response.GetKeyValue().Key,
		Value: response.GetKeyValue().Value,
	}, nil
}

func (p *ProxyService) PutKeyValue(key string, value string) (models.KeyValue, error) {
	req := &pb.KeyValue{
		Term:  1,
		Key:   key,
		Value: value,
	}

	response, err := p.requestForwarder.putKeyValue(*req)

	if err != nil {
		return models.KeyValue{}, err
	}

	return models.KeyValue{
		Key:   response.GetKeyValue().Key,
		Value: response.GetKeyValue().Value,
	}, nil
}

func (p *ProxyService) AddKeyValue(key string, value string) (models.KeyValue, error) {
	req := &pb.KeyValue{
		Term:  1,
		Key:   key,
		Value: value,
	}

	response, err := p.requestForwarder.addKeyValue(*req)

	if err != nil {
		return models.KeyValue{}, err
	}

	return models.KeyValue{
		Key:   response.GetKeyValue().Key,
		Value: response.GetKeyValue().Value,
	}, nil
}
