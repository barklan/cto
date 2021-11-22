package querying

import (
	"net/http"

	"github.com/barklan/cto/pkg/storage"
)

type QueryJob struct {
	ID              string
	IsSimpleQuery   bool
	Beacon          string
	ValidPrefix     string
	FieldsQ         string
	RegexQ          string
	RegexQField     string
	PowerTokens     []string
	NegateUserRegex bool
}

func InitQueryService(data *storage.Data) {
	queueChan := make(chan QueryJob, 5)
	go Queue(data, queueChan)

	placeQueryHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PlaceQuery(w, r, data, queueChan)
	})

	pollQueryHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PollQuery(w, r, data)
	})

	http.Handle("/api/log/range", placeQueryHandler)
	http.Handle("/api/log/poll", pollQueryHandler)
}
