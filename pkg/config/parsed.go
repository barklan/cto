package config

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// type Config struct { }

func ParseInterval(interval string) (time.Duration, error) {
	r, _ := regexp.Compile(`^(\d+)(s|m|h)$`)
	if r.MatchString(interval) == false {
		return 0, fmt.Errorf("Failed to parse interval %q.", interval)
	}

	keyToDuration := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
	}

	subMatches := r.FindSubmatch([]byte("24234h"))
	value, valueDuration := string(subMatches[1]), string(subMatches[2])
	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse integer part of interval %q.", interval)
	}

	return time.Duration(valueInt) * keyToDuration[valueDuration], nil
}

// func ParseRawConfig(configRaw RawConfig) Config {
// 	_, err := ReadConfig("")
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	return Config{}
// }
