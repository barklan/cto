package logserver

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/loginput"
	"github.com/barklan/cto/pkg/logserver/querying"
	"github.com/barklan/cto/pkg/logserver/types"
	"github.com/barklan/cto/pkg/storage"
)

type LogRecordReport struct {
	ProjectName string
}

type SessionData struct {
	Using            *uint64
	KnownErrorsMutex *sync.Mutex
	KnownErrors      []types.KnownError
}

func logOneRequest(
	projectName string,
	body []byte,
	data *storage.Data,
	sessDataMap map[string]*SessionData,
	reportChan chan LogRecordReport,
) {
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

		openOrEnterSession(data, sessDataMap, projectName)

		var wg sync.WaitGroup
		wg.Add(len(multiLog))
		for _, pick := range multiLog {
			go func(record RawLogRecord) {
				defer wg.Done()
				processLogRecord(data, projectName, record, sessDataMap, reportChan)
			}(pick)
		}
		wg.Wait()

		closeOrLeaveSession(data, sessDataMap, projectName)
	}(body)
}

func remove(s []types.KnownError, i int) []types.KnownError {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func processLogInputs(
	data *storage.Data,
	reqs <-chan loginput.LogRequest,
	sessDataMap map[string]*SessionData,
	reportChan chan LogRecordReport,
) {
	defer log.Panicln("log input processing stopped")
	for req := range reqs {
		logOneRequest(req.ProjectID, req.Body, data, sessDataMap, reportChan)
	}
}

func LogServerServe(data *storage.Data) {
	rand.Seed(time.Now().UnixNano())
	sessionDataMap := openSession(data)

	reportChan := make(chan LogRecordReport, 50)

	// FIXME this only works for one replica
	// 1. This should generate and save reports only for project that are resident on that node
	// 2. It should also periodically update list of projects
	projects := make([]string, 0)
	if err := data.R.Select(&projects, "select id from project"); err != nil {
		log.Println("no projects in database")
	}

	go func() {
		aggregateRecordsRecievedMap := map[string]int{}
		for _, projectName := range projects {
			aggregateRecordsRecievedMap[projectName] = 0
		}
		mu := sync.Mutex{}

		go func(m *sync.Mutex) {
			period := 1 * time.Minute
			ticker := time.NewTicker(period)
			for range ticker.C {
				for _, projectName := range projects {
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

	reqs := make(chan loginput.LogRequest, 10)
	go Subscriber(data, reqs)
	go processLogInputs(data, reqs, sessionDataMap, reportChan)

	logExactHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		querying.ServeLogExact(w, r, data)
	})
	http.Handle("/api/log/exact", logExactHandler)

	querying.InitQueryService(data)

	log.Panic(http.ListenAndServe(":8080", nil))
}
