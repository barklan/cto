package querying

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/barklan/cto/pkg/storage"
	"github.com/dgraph-io/badger/v3"
	"github.com/thedevsaddam/gojsonq/v2"
)

func jobFailed(data *storage.Data, key string, workerChan chan QueryJob) {
	data.SetObj(key, "error", 1*time.Hour)
}

// Worker should either set "timeout", "error", "none" or processed json data
// as a value for job.ID key.
func Worker(data *storage.Data, workerChan chan QueryJob) {
	for queryJob := range workerChan {

		log.Printf("recieved job %#v", queryJob)

		var userRegex *regexp.Regexp
		if queryJob.RegexQ != "" {
			var err error
			userRegex, err = regexp.Compile(queryJob.RegexQ)
			if err != nil {
				log.Printf("failed to compile user regexp; %v", err)
				jobFailed(data, queryJob.ID, workerChan)
				continue
			}
		}

		var requestedFields []string
		if queryJob.FieldsQ != "" {
			requestedFields = strings.Fields(queryJob.FieldsQ)
		}

		filteredValues := make([]map[string]interface{}, 0)
		err := data.DB.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.Reverse = true
			it := txn.NewIterator(opts)
			defer it.Close()
			prefix := []byte(queryJob.ValidPrefix)
			beacon := []byte(queryJob.Beacon)
			for it.Seek(beacon); it.ValidForPrefix(prefix); it.Next() {
				item := it.Item()
				// k := item.Key()
				var valCopy []byte
				err := item.Value(func(val []byte) error {
					valCopy = append([]byte{}, val...)
					processedVal, include, err := ProcessOneValue(
						valCopy,
						requestedFields,
						userRegex,
						queryJob.RegexQ,
						queryJob.RegexQField,
						queryJob.NegateUserRegex,
						queryJob.IsSimpleQuery,
					)
					if err != nil {
						log.Println("failed to process one value:", err)
					}
					if include {
						filteredValues = append(filteredValues, processedVal)
					}

					return nil
				})
				if err != nil {
					return err
				}

				if len(filteredValues) == 100 {
					break
				}
			}
			return nil
		})
		if err != nil {
			log.Println("failed to filter values:", err)
			jobFailed(data, queryJob.ID, workerChan)
			continue
		}

		if len(filteredValues) == 0 {
			data.SetObj(queryJob.ID, "none", 1*time.Hour)
		} else {
			data.SetObj(queryJob.ID, filteredValues, 1*time.Hour)
		}

		log.Println("worker finished job successfuly:", queryJob.ID)
	}
}

func ProcessOneValue(
	value []byte,
	requestedFields []string,
	userRegex *regexp.Regexp,
	regexQ string,
	regexQField string,
	negateUserRegex bool,
	isSimpleQuery bool,
) (map[string]interface{}, bool, error) {
	// TODO use fastjson here
	// - don't copy value before passing to function (check if this works first)
	// - separate requestedFields by dot before iterating
	// - introduce fastjson parser before iterating
	// - jqObj to fastjson value (p.Parse(string(value)))
	logrecord := map[string]interface{}{}

	err := json.Unmarshal(value, &logrecord)
	if err != nil {
		return nil, false, err
	}
	if isSimpleQuery {
		return logrecord, true, nil
	}

	filteredRawLogRecord := map[string]interface{}{}

	// TODO gojsonq is not maintained - write your own implementation
	jqObj := gojsonq.New().FromInterface(logrecord)

	if regexQ != "" {
		regexFieldValue := jqObj.Find(regexQField)
		if regexFieldValue != nil {
			switch regexFieldValue.(type) {
			case string:
				// TODO this can be done better
				if negateUserRegex {
					if !userRegex.MatchString(regexFieldValue.(string)) {
						if requestedFields == nil {
							return logrecord, true, nil
						}
					} else {
						return nil, false, nil
					}
				} else {
					if userRegex.MatchString(regexFieldValue.(string)) {
						if requestedFields == nil {
							return logrecord, true, nil
						}
					} else {
						return nil, false, nil
					}
				}
			default:
				return nil, false, nil
			}
		} else {
			return nil, false, nil
		}
	}

	// if we are here that means fields are not an empty string
	for _, requestedKey := range requestedFields {
		keyResult := jqObj.Reset().Find(requestedKey)
		if keyResult != nil {
			filteredRawLogRecord[requestedKey] = keyResult
		}
	}
	return filteredRawLogRecord, true, nil
}

func filterKeys(
	data *storage.Data,
	reg *regexp.Regexp,
	reverse bool,
	isSimpleQuery bool,
	cutOffHourOptim int64,
) [][]byte {
	var keys [][]byte
	var cutOffHourOptimStr string
	if cutOffHourOptim != -1 {
		cutOffHourOptimStr = strconv.FormatInt(cutOffHourOptim, 10)
		if cutOffHourOptim < 10 {
			cutOffHourOptimStr = "0" + cutOffHourOptimStr
		}
		cutOffHourOptimStr = fmt.Sprintf(" %s:", cutOffHourOptimStr)
	}
	err := data.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = reverse
		it := txn.NewIterator(opts)
		defer it.Close()

		i := 0
		for it.Rewind(); it.Valid(); it.Next() {
			i++
			item := it.Item()
			k := item.Key()
			if reg.Match(k) == true {
				keyCopy := append([]byte{}, k...)
				keys = append(keys, keyCopy)
				if isSimpleQuery && len(keys) == 100 {
					break
				}
			} else {
				if cutOffHourOptim != -1 && i%40 == 0 { // 40 should be enough
					if strings.Contains(string(k), cutOffHourOptimStr) {
						break
					}
				}
			}

		}
		return nil
	})
	if err != nil {
		log.Println("ERROR iterating")
	}
	return keys
}

func simpleQuery(data *storage.Data, keys [][]byte) []byte {
	values := make([]byte, 0)
	newline := []byte(",")

	values = append(values, []byte("[")...)
	for i, key := range keys {
		value := data.GetLogRaw(key)
		values = append(values, value...)
		if i != len(keys)-1 {
			values = append(values, newline...)
		}
	}
	values = append(values, []byte("]")...)
	return values
}
