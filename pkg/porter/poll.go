package porter

import (
	"net/http"
)

func pollLogRange(base *Base, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	rawQuery := r.URL.Query()

	tokenQ := rawQuery.Get("token")
	_, _, statusCode, ok := authorize(base, tokenQ)
	if !ok {
		w.WriteHeader(statusCode)
		return
	}

	qID := rawQuery.Get("qid")

	query, ok, err := base.Cache.Get(qID)
	if err != nil {
		http.Error(w, "Failed to get query from cache.", 500)
		return
	}
	if !ok {
		http.Error(w, "Query request with this id not found.", 404)
		return
	}

	w.WriteHeader(200)
	_, _ = w.Write(query)
}
