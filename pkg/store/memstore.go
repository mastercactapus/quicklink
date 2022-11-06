package store

import (
	"context"
	"sync"
)

type MemStore struct {
	m map[string]string

	mx sync.RWMutex
}

func NewMemStore() *MemStore {
	return &MemStore{
		m: make(map[string]string),
	}
}

func (s *MemStore) Get(ctx context.Context, key string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.m[key], nil
}

func (s *MemStore) Set(ctx context.Context, key, value string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	if value == "" {
		delete(s.m, key)
		return nil
	}

	s.m[key] = value
	return nil
}

func (s *MemStore) Scanner(ctx context.Context) (Scanner, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	var keys []entry
	for k := range s.m {
		keys = append(keys, entry{k, s.m[k]})
	}

	return &entryScanner{
		results: keys,
	}, nil
}
