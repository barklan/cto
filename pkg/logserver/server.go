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
	"github.com/barklan/cto/pkg/postgres/models"
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

func authorizeRequest(data *storage.Data, r *http.Request) (string, bool) {
	projectName, password, ok := r.BasicAuth()
	if !ok {
		log.Println("error parsing basic auth")
		return "", false
	}

	project := models.Project{}
	if err := data.R.Get(&project, "select * from project where id = $1", projectName); err != nil {
		log.Println("project not found from basic auth")
		return "", false
	}

	if subtle.ConstantTimeCompare([]byte(password), []byte(project.SecretKey)) == 1 {
		return project.ID, true
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

func LogServerServe(data *storage.Data) {
	rand.Seed(time.Now().UnixNano())
	sessionDataMap := openSession(data)

	reportChan := make(chan LogRecordReport, 50)

	projects := make([]string, 0)
	if err := data.R.Select(&projects, "select id from project"); err != nil {
		log.Println("no projects found when opening logserver session")
	}

	// Save new reports every 5 minutes
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
