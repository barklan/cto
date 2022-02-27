package storage

import "testing"

func TestSetGet(t *testing.T) {
	db := OpenDB("/tmp/badger_test", "/main")
	defer db.Close()
	key := "foo"
	value := "bar"

	Set(db, key, []byte(value))
	valueActual := Get(db, key)

	if string(valueActual) != value {
		t.Errorf("got %v want %v", valueActual, value)
	}
}
