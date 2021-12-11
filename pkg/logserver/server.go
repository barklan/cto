package logserver

import (
	"crypto/subtle"
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
	"github.com/barklan/cto/pkg/logserver/types"
	"github.com/barklan/cto/pkg/storage"
)

var LogServerSessionDataMap map[string]*SessionData

type LogRecordReport struct {
	ProjectName string
}

type SessionData struct {
	KnownErrorsMutex sync.Mutex
	KnownErrors      []types.KnownError
}

func authorizeRequest(data *storage.Data, r *http.Request) (string, bool) {
	_, password, ok := r.BasicAuth()
	if !ok {
		log.Println("error parsing basic auth")
		return "", false
	}

	for projectName, config := range data.Config.P {
		if subtle.ConstantTimeCompare([]byte(password), []byte(config.SecretKey)) == 1 {
			return projectName, true
		}
	}

	return "", false
}

func logOneRequest(
	w http.ResponseWriter,
	r *http.Request,
	data *storage.Data,
	sessDataMap map[string]*SessionData,
	reportChan chan LogRecordReport,
) {
	projectName, ok := authorizeRequest(data, r)
	if !ok {
		log.Println("recieved unauthorized request")
		w.WriteHeader(401)
		return
	}

	log.Println(fmt.Sprintf("recieved log dump for project %q", projectName))

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

		var wg sync.WaitGroup
		wg.Add(len(multiLog))

		for _, pick := range multiLog {
			go func(record RawLogRecord) {
				defer wg.Done()
				processLogRecord(data, record, sessDataMap, reportChan)
			}(pick)
		}

		wg.Wait()
	}(body)
}

func remove(s []types.KnownError, i int) []types.KnownError {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func LogServerServe(data *storage.Data) {
	rand.Seed(time.Now().UnixNano())
	sessionDataMap := openSession(data)
	LogServerSessionDataMap = sessionDataMap

	// TODO maybe use this for websocket watch mode?
	reportChan := make(chan LogRecordReport, 50) // TODO maybe not enough

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
					report := types.PeriodicReport{
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
	err := http.ListenAndServe(portString, nil)
	if err != nil {
		_, _ = data.CSendSync("Logserver errored.")
		log.Panic(err)
	}
}
