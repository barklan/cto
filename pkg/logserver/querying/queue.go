package querying

import (
	"container/list"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/barklan/cto/pkg/storage"
)

func Queue(data *storage.Data, queueChan chan QueryJob) {
	log.Println("queue starting")

	jobsQueue := list.New()
	var mx sync.Mutex

	workerChan := make(chan QueryJob)

	go Worker(data, workerChan)
	go Worker(data, workerChan)

	go func(jobsQueue *list.List, m *sync.Mutex) {
		for {
			m.Lock()
			if jobsQueue.Len() > 0 {
				job := jobsQueue.Front()
				m.Unlock()
				jobValue := job.Value.(QueryJob)
				workerChan <- jobValue
				data.SetObj(jobValue.ID, "processing", 1*time.Hour)
				m.Lock()
				jobsQueue.Remove(job)
				m.Unlock()
			} else {
				m.Unlock()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}(jobsQueue, &mx)

	for requestedJob := range queueChan {
		log.Println("queue recieved new job")

		data.SetObj(requestedJob.ID, "queued", 1*time.Hour)

		mx.Lock()
		jobsQueue.PushBack(requestedJob)
		mx.Unlock()
	}
}
