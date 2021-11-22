package logserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/logserver/querying"
	"github.com/barklan/cto/pkg/storage"
)

type KnownError struct {
	OriginBadgerKey string
	LogStr          string
	Counter         uint64 // should use atomic operations on this one
	LastSeen        time.Time
}

type LogRecordReport struct {
	ProjectName string
}

type SessionData struct {
	KnownErrorsMutex sync.Mutex
	KnownErrors      []KnownError
}

func logOneRequest(
	w http.ResponseWriter,
	r *http.Request,
	data *storage.Data,
	sessDataMap map[string]*SessionData,
	reportChan chan LogRecordReport,
) {
	log.Println("recieved log dump")
	body, _ := io.ReadAll(r.Body)

	go func([]byte) {
		multiLog := make([]RawLogRecord, 1)
		err := json.Unmarshal(body, &multiLog)
		if err != nil {
			log.Println("Failed to unmarshal log request.")
		}

		recordsRecieved := len(multiLog)
		log.Printf("unmarshaled log dump containing %d records", recordsRecieved)
		if recordsRecieved == 0 {
			log.Println("no records to unmarshal")
			return
		}

		// subset := GetSubset(multiLog, 5)
		var wg sync.WaitGroup
		wg.Add(len(multiLog))
		// randomIndex := rand.Intn(recordsRecieved)
		// pick := multiLog[randomIndex]
		for _, pick := range multiLog {
			go func(record RawLogRecord) {
				defer wg.Done()
				processLogRecord(data, record, sessDataMap, reportChan)
			}(pick)
		}

		wg.Wait()
	}(body)
}

func remove(s []KnownError, i int) []KnownError {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type PeriodicReport struct {
	Period   time.Duration
	Recieved int
}

func LogServerServe(data *storage.Data) {
	rand.Seed(time.Now().UnixNano())
	sessionDataMap := openSession(data)

	// TODO maybe use this for websocket watch mode?
	reportChan := make(chan LogRecordReport, 20) // TODO maybe not enough

	// Save new reports every 5 minutes
	go func() {
		aggregateRecordsRecievedMap := map[string]int{}
		for projectName := range data.Config.P {
			aggregateRecordsRecievedMap[projectName] = 0
		}
		mu := sync.Mutex{}

		go func(m *sync.Mutex) {
			period := 5 * time.Minute
			ticker := time.NewTicker(period)
			for range ticker.C {
				for projectName := range data.Config.P {
					report := PeriodicReport{
						Period:   period,
						Recieved: aggregateRecordsRecievedMap[projectName],
					}
					data.SetObj(fmt.Sprintf("periodicLogReport-%s", projectName), report, period)

					m.Lock()
					aggregateRecordsRecievedMap[projectName] = 0
					m.Unlock()
				}
			}
		}(&mu)

		for report := range reportChan {
			mu.Lock()
			aggregateRecordsRecievedMap[report.ProjectName]++
			mu.Unlock()

		}
	}()

	// TODO should use basic auth
	logInputHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logOneRequest(w, r, data, sessionDataMap, reportChan)
	})
	http.Handle("/api/log/input", logInputHandler)

	logExactHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		querying.ServeLogExact(w, r, data)
	})
	http.Handle("/api/log/exact", logExactHandler)

	querying.InitQueryService(data)

	var portString string
	if v, ok := os.LookupEnv("CTO_LOCAL_ENV"); ok {
		if v == "true" {
			portString = ":8080"
		}
	} else {
		portString = ":8080"
	}
	http.ListenAndServe(portString, nil)
}
