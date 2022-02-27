package types

import (
	"encoding/json"
	"time"
)

type KnownError struct {
	LastSeen        time.Time `json:"last_seen,omitempty"`
	Hostname        string    `json:"hostname,omitempty"`
	Service         string    `json:"service,omitempty"`
	OriginBadgerKey string    `json:"origin_badger_key,omitempty"`
	LogStr          string    `json:"log_str,omitempty"`
	Counter         uint64    `json:"counter,omitempty"` // should use atomic operations on this one
}

func (e KnownError) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}
