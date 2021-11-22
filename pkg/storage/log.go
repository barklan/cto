package storage

import (
	"encoding/json"
	"log"
	"reflect"
	"time"
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
		SetWithTTL(d.LogDB, key, byteObj, ttl)
	} else {
		log.Panic("WTF. Logs should have ttl")
	}
}

// You need to unmarshal it yourself.
func (d *Data) GetLog(key string) []byte {
	return Get(d.LogDB, key)
}

func (d *Data) GetLogRaw(key []byte) []byte {
	return GetRaw(d.LogDB, key)
}
