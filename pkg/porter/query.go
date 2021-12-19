package porter

import (
	"net/http"

	"github.com/barklan/cto/pkg/bot"
	log "github.com/sirupsen/logrus"
)

// This needs to be secured
func serveLogExact(base *Base, s *bot.Sylon, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	query := r.URL.Query()

	badgerKey := query.Get("key")

	value, ok, err := base.Cache.Get(badgerKey)
	if err != nil {
		s.CSend("error when getting value from cache")
		http.Error(w, "error when getting value from cache", 500)
		return
	}
	if !ok {
		http.Error(w, "value for that key is not in cache.", 404)
		return
	}

	log.Info(string(value))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(value)
}
