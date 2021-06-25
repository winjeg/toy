package store

import "time"

type BasicStore interface {
	Set(k, v []byte) error
	Get(k []byte) ([]byte, error)
	SetEx(k, v []byte, exp time.Duration) error
	Del(k []byte) error
}

type Seq struct {
	Max  uint64
	Next func() uint64
}

var (
	store BasicStore = &badgerStore{getBadger()}
)

func Set(k, v string) error {
	return store.Set([]byte(k), []byte(v))
}
func Get(k string) ([]byte, error) {
	return store.Get([]byte(k))
}

func SetEx(k, v string, exp time.Duration) error {
	return store.SetEx([]byte(k), []byte(v), exp)
}

func Del(k string) error {
	return store.Del([]byte(k))
}
