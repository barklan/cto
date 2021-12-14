package querying

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/barklan/cto/pkg/storage"
)

// This needs to be secured
func ServeLogExact(w http.ResponseWriter, r *http.Request, data *storage.Data) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query()

	badgerKey := query.Get("key")

	value := data.GetLog(badgerKey)
	if string(value) == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO just return bytes directly from badger.
	var valueDec interface{}
	err := json.Unmarshal(value, &valueDec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("FAILED to unmarshal")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(value)
}
