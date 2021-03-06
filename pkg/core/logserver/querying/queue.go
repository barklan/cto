package querying

import (
	"container/list"
	"log"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/core/storage"
	"github.com/barklan/cto/pkg/porter"
)

func Queue(data *storage.Data, queueChan chan QueryJob) {
	data.Log.Info("queue starting")

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
				SetMsgInCache(data, jobValue.ID, porter.QWorking, "Worker started processing.")
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

		SetMsgInCache(data, requestedJob.ID, porter.QWorking, "Query was queued in core node.")

		mx.Lock()
		jobsQueue.PushBack(requestedJob)
		mx.Unlock()
	}
}
