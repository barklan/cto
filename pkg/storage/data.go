package storage

import (
	"encoding/json"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/barklan/cto/pkg/caching"
	"github.com/dgraph-io/badger/v3"
	"github.com/jmoiron/sqlx"
	tb "gopkg.in/tucnak/telebot.v2"
)

var variableKeySymbol = "$"

type Data struct {
	Chat      *tb.Chat
	DB        *badger.DB
	Config    *Config
	MediaPath string
	R         *sqlx.DB
	Cache     caching.Cache
}

func (d *Data) SetVar(projectName, key string, obj interface{}, ttl time.Duration) {
	var byteObj []byte
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.String {
		byteObj = []byte(obj.(string))
	} else {
		var err error
		byteObj, err = json.Marshal(obj)
		if err != nil {
			log.Panic(err)
		}
	}

	varKey := projectName + variableKeySymbol + key
	if ttl > 0 {
		SetWithTTL(d.DB, varKey, byteObj, ttl)
	} else {
		Set(d.DB, varKey, byteObj)
	}
}

func (d *Data) GetVar(projectName, key string) []byte {
	varKey := projectName + variableKeySymbol + key
	return Get(d.DB, varKey)
}

func (d *Data) VarExists(projectName, key string) bool {
	varKey := projectName + variableKeySymbol + key
	err := d.DB.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(varKey))
		return err
	})
	switch err {
	case badger.ErrKeyNotFound:
		return false
	case nil:
		return true
	default:
		log.Panicln("panicing in VarExists function", err)
		return false
	}
}

func (d *Data) DeleteVar(projectName, key string) {
	varKey := projectName + variableKeySymbol + key
	err := d.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(varKey))
	})
	if err != nil {
		log.Panicln("panicing in DeleteVar function", err)
	}
}
