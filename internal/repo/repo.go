package repo

import (
	"errors"
	"raft-based-kv/internal/models"
)

type Repo struct {
	in_memory_db map[string]string
}

func NewRepository() *Repo {
	newRepo := Repo{
		in_memory_db: make(map[string]string),
	}
	return &newRepo
}

func (s *Repo) GetKeyValue(key string) (models.KeyValue, error) {

	value := s.in_memory_db[key]

	if value == "" {
		err := errors.New("Could not find a Key Value definition")
		return models.KeyValue{}, err
	}

	return models.KeyValue{
		Key:   key,
		Value: value,
	}, nil
}

func (s *Repo) AddKeyValue(key string, value string) (models.KeyValue, error) {
	//First Check if they Key value pair does not already exist
	if s.in_memory_db[key] != "" {
		err := errors.New("The Key is aready assigned a value use the PUT endpoint to update key value")
		return models.KeyValue{}, err
	}
	s.in_memory_db[key] = value
	return models.KeyValue{
		Key:   key,
		Value: s.in_memory_db[key],
	}, nil
}

func (s *Repo) DeleteKeyValue(key string) error {
	value := s.in_memory_db[key]

	if value == "" {
		err := errors.New("Could not find a key value definition")
		return err
	}
	delete(s.in_memory_db, key)
	return nil
}

func (s *Repo) UpdateKeyValue(key string, value string) (models.KeyValue, error) {
	checkValue := s.in_memory_db[key]

	if checkValue == "" {
		err := errors.New("Could not find a key value definition")
		return models.KeyValue{}, err
	}

	s.in_memory_db[key] = value
	return models.KeyValue{
		Key:   key,
		Value: s.in_memory_db[key],
	}, nil
}
