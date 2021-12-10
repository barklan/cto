package types

import "time"

type PeriodicReport struct {
	Period   time.Duration
	Recieved int
}
