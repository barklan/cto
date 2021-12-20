package porter

import (
	"net/http"
)

func pollLogRange(base *Base, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	rawQuery := r.URL.Query()

	tokenQ := rawQuery.Get("token")
	_, statusCode, ok := authorize(base, tokenQ)
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

	// qResp := QResp{}
	// if err := json.Unmarshal(query, qResp); err != nil {
	// 	w.Write([]byte("Failed to unmarshal query request from cache."))
	// 	w.WriteHeader(500)
	// 	return
	// }

	w.WriteHeader(200)
	_, _ = w.Write(query)

	// status := qResp.Status

	// switch status {
	// case QWorking:
	// 	w.WriteHeader(statusCode int)
	// }

	// queryResult := base.Get(qID)
	// queryResultStr := string(queryResult)
	// switch queryResultStr {
	// case "":
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// case "error":
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// case "processing":
	// 	w.WriteHeader(http.StatusOK) // 200 indicates that frontend should continue polling
	// 	return
	// case "queued":
	// 	w.WriteHeader(http.StatusAccepted)
	// 	return
	// case "none":
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// default:
	// 	w.WriteHeader(http.StatusCreated) // 201 indicates that the data is available
	// 	w.Write(queryResult)
	// 	return
	// }
}
