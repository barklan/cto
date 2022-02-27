package querying

import (
	"log"

	"github.com/barklan/cto/pkg/core/storage"
	"github.com/barklan/cto/pkg/porter"
)

type QueryJob struct {
	ID              string
	IsSimpleQuery   bool
	Beacon          string
	ValidPrefix     string
	NorthStar       string
	FieldsQ         string
	RegexQ          string
	RegexQField     string
	NegateUserRegex bool
}

func InitQueryService(data *storage.Data) {
	defer log.Panic("query block exited")

	queueChan := make(chan QueryJob, 5)
	reqs := make(chan porter.QueryRequest, 10)

	go Queue(data, queueChan)
	go Subscriber(data, reqs)

	for query := range reqs {
		PlaceQuery(query, data, queueChan)
	}
}
