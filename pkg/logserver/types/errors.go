package types

import "time"

type KnownError struct {
	LastSeen        time.Time
	Hostname        string
	Service         string
	OriginBadgerKey string
	LogStr          string
	Counter         uint64 // should use atomic operations on this one
}
