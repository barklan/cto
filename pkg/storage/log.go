package storage

import (
	"encoding/json"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
)

// You do not need to marshal anything!
func (d *Data) SetLog(key string, obj interface{}, ttl time.Duration) {
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

	if ttl > 0 {
		SetWithTTL(d.DB, key, byteObj, ttl)
	} else {
		log.Panic("WTF. Logs should have ttl")
	}
}

// You need to unmarshal it yourself.
func (d *Data) GetLog(key string) []byte {
	return Get(d.DB, key)
}

func (d *Data) GetLogRaw(key []byte) []byte {
	return GetRaw(d.DB, key)
}
