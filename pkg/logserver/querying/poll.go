package querying

import (
	"net/http"

	"github.com/barklan/cto/pkg/storage"
)

func PollQuery(w http.ResponseWriter, r *http.Request, data *storage.Data) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	rawQuery := r.URL.Query()

	tokenQ := rawQuery.Get("token")
	_, statusCode, ok := authorize(data, tokenQ)
	if !ok {
		w.WriteHeader(statusCode)
		return
	}

	qID := rawQuery.Get("qid")
	queryResult := data.Get(qID)
	queryResultStr := string(queryResult)
	switch queryResultStr {
	case "":
		w.WriteHeader(http.StatusBadRequest)
		return
	case "error":
		w.WriteHeader(http.StatusInternalServerError)
		return
	case "processing":
		w.WriteHeader(http.StatusOK) // 200 indicates that frontend should continue polling
		return
	case "queued":
		w.WriteHeader(http.StatusAccepted)
		return
	case "none":
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		w.WriteHeader(http.StatusCreated) // 201 indicates that the data is available
		w.Write(queryResult)
		return
	}
}
