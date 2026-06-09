package main

import (
	"sync"
)

type Store struct {
	kv_map map[string]string
	mtx    sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		kv_map: make(map[string]string),
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	value, ok := s.kv_map[key]
	return value, ok
}

func (s *Store) Set(key, value string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.kv_map[key]
	s.kv_map[key] = value
	return !ok
}

func (s *Store) Delete(key string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.kv_map[key]
	if ok {
		delete(s.kv_map, key)
	}
	return ok
}

func (s *Store) Keys() []string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	keys := make([]string, 0, len(s.kv_map))
	for k := range s.kv_map {
		keys = append(keys, k)
	}
	return keys
}

func (s *Store) Count() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return len(s.kv_map)
}
