package storage

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func SetData() *Data {
	db := OpenDB("/tmp/badger_test", "/main")
	data := &Data{DB: db}
	return data
}

func RandomKey() string {
	rand.Seed(time.Now().UnixNano())
	key := fmt.Sprintf("foo%d", rand.Intn(10000))
	return key
}

func TestSetObj(t *testing.T) {
	t.Run("Set string obj", func(t *testing.T) {
		data := SetData()
		defer data.DB.Close()

		key := RandomKey()
		value := "bar"
		data.SetObj(key, value, 1*time.Minute)

		got := data.GetStr(key)
		if got != value {
			t.Errorf("got %v want %v", got, value)
		}
	})
	t.Run("Set string with ttl", func(t *testing.T) {
		data := SetData()
		defer data.DB.Close()

		key := RandomKey()
		value := "bar2"
		data.SetObj(key, value, 10*time.Millisecond)

		time.Sleep(12 * time.Millisecond)
		got := data.GetStr(key)

		if got != "" {
			t.Errorf("Expected to get empty string on expired key.")
		}
	})
}
