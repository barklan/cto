package storage

import (
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger/v3"
)

func OpenDB(customPath string, dbPath string) *badger.DB {
	var badgerPath string
	if customPath != "" {
		badgerPath = customPath
	} else if v, ok := os.LookupEnv("CTO_DATA_PATH"); ok {
		badgerPath = v + dbPath
	} else {
		badgerPath = "/app/data" + dbPath
	}
	log.Printf("setting badger path %q", badgerPath)
	badgerOptions := badger.DefaultOptions(badgerPath)
	// badgerOptions.SyncWrites = true
	db, err := badger.Open(badgerOptions)
	if err != nil {
		log.Panic(err)
	}
	return db
}

func Set(db *badger.DB, key string, value []byte) {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), value)
		return err
	})
	if err != nil {
		log.Panicf("Failed to set kv to badger. Error: %v", err)
	}
}

func SetWithTTL(db *badger.DB, key string, value []byte, ttl time.Duration) {
	err := db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), value).WithTTL(ttl)
		err := txn.SetEntry(e)
		return err
	})
	if err != nil {
		log.Panicf("Failed to set kv with ttl to badger. Error: %v", err)
	}
}

func Get(db *badger.DB, key string) []byte {
	rawKey := []byte(key)
	return GetRaw(db, rawKey)
}

func GetRaw(db *badger.DB, key []byte) []byte {
	var valCopy []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		return []byte("")
	}

	return valCopy
}
