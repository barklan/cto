package types

import "time"

type KnownError struct {
	OriginBadgerKey string
	LogStr          string
	Counter         uint64 // should use atomic operations on this one
	LastSeen        time.Time
}
