package store

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"

	"fmt"
	"log"
	"sync"
	"time"
)

var (
	once     = sync.Once{}
	badgerDb *badger.DB
)

func defaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	return logger
}

func getBadger() *badger.DB {
	if badgerDb != nil {
		return badgerDb
	}
	once.Do(func() {
		opts := badger.DefaultOptions("E:/Desktop/badger")
		opts.Logger = defaultLogger()
		db, err := badger.Open(opts)
		if err != nil {
			log.Fatal(err)
			return
		}
		badgerDb = db
	})
	return badgerDb
}

func CleanUp() {
	if badgerDb == nil {
		return
	}
	err := badgerDb.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
}

type badgerStore struct {
	DB *badger.DB
}

func (bs *badgerStore) Set(k, v []byte) error {
	err := bs.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(k, v)
		return err
	})
	return err
}

func (bs *badgerStore) Get(k []byte) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	var result []byte
	err := bs.DB.View(func(txn *badger.Txn) error {
		item, _ := txn.Get(k)
		err := item.Value(func(val []byte) error {
			result = append([]byte{}, val...)
			return nil
		})
		return err
	})
	return result, err
}

func (bs *badgerStore) SetEx(k, v []byte, exp time.Duration) error {
	err := bs.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(k, v).WithTTL(exp)
		err := txn.SetEntry(e)
		return err
	})
	return err
}

func (bs *badgerStore) Seq(k []byte, max uint64) func() (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Wrong key: ", r)
			return
		}
	}()
	seq, err := bs.DB.GetSequence(k, max)
	if err != nil {
		return nil
	}
	defer func(seq *badger.Sequence) {
		err := seq.Release()
		if err != nil {
			log.Println(err.Error())
		}
	}(seq)
	return func() (uint64, error) {
		return seq.Next()
	}
}

func (bs *badgerStore) Del(k []byte) error {
	err := bs.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete(k)
		return err
	})
	return err
}
