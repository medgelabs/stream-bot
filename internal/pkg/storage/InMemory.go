package storage

import "fmt"

type InMemory struct {
	store map[string]interface{}
}

func NewInMemory() *InMemory {
	return &InMemory{
		store: make(map[string]interface{}),
	}
}

func (s InMemory) put(key string, value interface{}) error {
	s.store[key] = value
	return nil
}

func (s InMemory) get(key string) (interface{}, error) {
	value, ok := s.store[key]
	if !ok {
		return nil, fmt.Errorf("could not get key")
	}

	return value, nil
}
