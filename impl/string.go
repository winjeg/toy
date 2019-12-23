package impl

type strStore struct {
	store map[string][]byte
}

func (s *strStore) Set(k string, v []byte) error {
	s.store[k] = v
	return nil
}
func (s *strStore) SetEx(k string , v []byte, exp int) error {
	s.store[k] = v
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
