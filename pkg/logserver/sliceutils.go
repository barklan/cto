package logserver

import (
	"math/rand"
	"time"
)

func GetSubset(array []RawLogRecord, divisor int) []RawLogRecord {
	var count int
	if len(array) < divisor {
		count = 1
	} else {
		count = len(array) / divisor
	}
	return getRandomElements(array, count)
}

func getRandomElements(array []RawLogRecord, count int) []RawLogRecord {
	result := make([]RawLogRecord, 0)
	existingIndexes := make(map[int]struct{}, 0)
	randomElementsCount := count

	for i := 0; i < randomElementsCount; i++ {
		randomIndex := randomIndex(len(array), existingIndexes)
		result = append(result, array[randomIndex])
	}

	return result
}

func randomIndex(size int, existingIndexes map[int]struct{}) int {
	rand.Seed(time.Now().UnixNano())

	for {
		randomIndex := rand.Intn(size)

		_, exists := existingIndexes[randomIndex]
		if !exists {
			existingIndexes[randomIndex] = struct{}{}
			return randomIndex
		}
	}
}
