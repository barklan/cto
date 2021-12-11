package types

import "time"

type KnownError struct {
	Hostname        string
	Service         string
	OriginBadgerKey string
	LogStr          string
	Counter         uint64 // should use atomic operations on this one
	LastSeen        time.Time
}
