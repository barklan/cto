package querying

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/barklan/cto/pkg/core/storage"
	"github.com/barklan/cto/pkg/porter"
	"github.com/dgraph-io/badger/v3"
	"github.com/thedevsaddam/gojsonq/v2"
)

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
				SetMsgInCache(data, queryJob.ID, porter.QFailed, "Failed to compile user regexp in worker.")
				continue
			}
		}

		var requestedFields []string
		if queryJob.FieldsQ != "" {
			requestedFields = strings.Fields(queryJob.FieldsQ)
		}

		filteredValues := make([]map[string]interface{}, 0)
		tooMuch := false
		err := data.DB.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.Reverse = true
			it := txn.NewIterator(opts)
			defer it.Close()

			northStarKey := []byte{}
			if queryJob.NorthStar != "" {
				it.Seek([]byte(queryJob.NorthStar))
				if it.Valid() {
					northStarKey = it.Item().KeyCopy(northStarKey)
					log.Println("northStarKey:", string(northStarKey))
				} else {
					log.Printf("Iteration is not valid")
					northStarKey = []byte{}
				}
				it.Rewind()
			}

			prefix := []byte(queryJob.ValidPrefix)
			beacon := []byte(queryJob.Beacon)
			for it.Seek(beacon); it.ValidForPrefix(prefix); it.Next() {
				item := it.Item()
				if bytes.Equal(item.Key(), northStarKey) {
					break
				}
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
					tooMuch = true
					break
				}
			}
			return nil
		})
		if err != nil {
			log.Println("failed to filter values:", err)
			SetMsgInCache(data, queryJob.ID, porter.QFailed, "Internal error. Failed to filter values.")
			continue
		}

		if len(filteredValues) == 0 {
			SetMsgInCache(data, queryJob.ID, porter.QFailed, "No results found.")
		} else {
			var msg string
			if tooMuch {
				msg = "Too many matching events found. Only the last 100 are shown."
			} else {
				msg = fmt.Sprintf("%d matching events found.", len(filteredValues))
			}
			SetResultInCache(data, queryJob.ID, msg, filteredValues)
		}

		data.Log.Info("worker finished job successfuly", zap.String("qid", queryJob.ID))
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
		data.Log.Error("error iterating", zap.Error(err))
	}
	return keys
}
