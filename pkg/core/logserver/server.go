package logserver

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/barklan/cto/pkg/core/logserver/querying"
	"github.com/barklan/cto/pkg/core/logserver/types"
	"github.com/barklan/cto/pkg/core/storage"
	"github.com/barklan/cto/pkg/loginput"
	"github.com/barklan/cto/pkg/vars"
)

type LogRecordReport struct {
	ProjectName string
}

type SessionData struct {
	Using       *uint64
	Mutex       *sync.Mutex
	KnownErrors map[string]types.KnownError
}

func logOneRequest(
	projectName string,
	body []byte,
	data *storage.Data,
	sessDataMap map[string]*SessionData,
	reportChan chan LogRecordReport,
) {
	multiLog := make([]RawLogRecord, 1)
	err := json.Unmarshal(body, &multiLog)
	if err != nil {
		data.Log.Warn("failed to unmarshal log request", zap.String("project", projectName))
	}

	recordsRecieved := len(multiLog)
	data.Log.Info("unmarshaled log dump", zap.String("project", projectName), zap.Int("records", recordsRecieved))
	if recordsRecieved == 0 {
		data.Log.Warn("no records to unmarshal", zap.String("project", projectName))
		return
	}

	sessData := openOrEnterSession(data, sessDataMap, projectName)

	subSet := GetSubset(multiLog, 20)

	var wg sync.WaitGroup
	wg.Add(len(subSet))
	for _, pick := range subSet {
		go func(record RawLogRecord) {
			defer wg.Done()
			processLogRecord(data, projectName, record, sessData, reportChan)
		}(pick)
	}
	wg.Wait()

	closeOrLeaveSession(data, sessDataMap, projectName)
}

func processLogInputs(
	data *storage.Data,
	reqs <-chan loginput.LogRequest,
	sessDataMap map[string]*SessionData,
	reportChan chan LogRecordReport,
) {
	defer data.Log.Panic("log input processing stopped")
	for req := range reqs {
		bodyCopy := make([]byte, len(req.Body))
		copy(bodyCopy, req.Body)
		go logOneRequest(req.ProjectID, bodyCopy, data, sessDataMap, reportChan)
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
					data.SetVar(projectName, vars.PeriodicLogReport, report, period)

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

	querying.InitQueryService(data)
}
