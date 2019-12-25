package impl

import "sync"

type strStore struct {
	lock  sync.Mutex
	store map[string][]byte
}

func (s *strStore) Set(k string, v []byte) error {
	s.lock.Lock()
	s.store[k] = v
	s.lock.Unlock()
	return nil
}

func (s *strStore) SetEx(k string, v []byte, exp int) error {
	s.lock.Lock()
	s.store[k] = v
	s.lock.Unlock()
	return nil
}

func (s *strStore) Get(k string) []byte {
	return s.store[k]
}

func NewStrStore() *strStore {
	return &strStore{
		store: make(map[string][]byte, 128),
	}
}
