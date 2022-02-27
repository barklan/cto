package logserver

import (
	"reflect"
	"testing"
)

func TestGetSubset(t *testing.T) {
	set := []RawLogRecord{
		{"one": 1},
		{"two": 2},
		{"three": 3},
		{"four": 4},
		{"five": 5},
		{"six": 6},
		{"seven": 7},
		{"eight": 8},
		{"nine": 9},
		{"ten": 10},
		{"eleven": 11},
		{"twelve": 12},
		{"thirteen": 13},
		{"fourteen": 14},
		{"fifteen": 15},
		{"sixteen": 16},
		{"seventeen": 17},
		{"eightteen": 18},
		{"nineteen": 19},
		{"twenty": 20},
	}

	for i := 0; i < 100; i++ {
		got := GetSubset(set, 5)
		for j, v := range got {
			for k := j + 1; k < len(got); k++ {
				if reflect.DeepEqual(v, got[k]) {
					t.Errorf("Fuck")
				}
			}
		}
	}
}
