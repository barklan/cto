package porter

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/barklan/cto/pkg/bot"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

type QStatus int

const (
	QWorking QStatus = iota
	QDone
	QFailed
)

type QResp struct {
	// TODO don't use magic strings
	// Status should be one of: working, failed, done
	Status QStatus                  `json:"status,omitempty"`
	Msg    string                   `json:"msg,omitempty"`
	Result []map[string]interface{} `json:"result,omitempty"`
}

// This needs to be secured
func serveLogExact(
	base *Base,
	s *bot.Sylon,
	w http.ResponseWriter,
	r *http.Request,
) {
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

func SetQRespInCache(base *Base, requestId string, status QStatus, msg string) {
	key := requestId

	valJson, err := json.Marshal(QResp{Msg: msg, Status: status})
	if err != nil {
		log.Panicln("failed to marshal meta message for requested query", err)
	}

	if err := base.Cache.Set(key, valJson, 1*time.Minute); err != nil {
		log.Error("failed to set qmeta in cache", err)
	}
}

func serveLogRange(
	base *Base,
	s *bot.Sylon,
	queries chan<- QueryRequestWrap,
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	rawQuery := r.URL.Query()

	tokenQ := rawQuery.Get("token")
	projectName, statusCode, ok := authorize(base, tokenQ)
	if !ok {
		w.WriteHeader(statusCode)
		return
	}

	query := rawQuery.Get("query")
	fieldsQ := rawQuery.Get("fields")
	regexQRaw := rawQuery.Get("regex")

	uid4, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Failed to generate request uuid.", 500)
		return
	}
	u4 := uid4.String()

	qreq := QueryRequest{
		RequestID: u4,
		ProjectID: projectName,
		QueryText: query,
		Fields:    fieldsQ,
		Regex:     regexQRaw,
	}

	qreqJson, err := json.Marshal(qreq)
	if err != nil {
		http.Error(w, "Failed to serialize request.", 500)
	}

	wrapped := QueryRequestWrap{
		ProjectID: projectName,
		Json:      qreqJson,
	}

	queries <- wrapped

	respMap := map[string]string{"qid": u4}
	resp, err := json.Marshal(respMap)
	if err != nil {
		http.Error(w, "failed to marshal response:", 500)
		return
	}

	SetQRespInCache(base, u4, QWorking, "Query request was accepted.")

	w.WriteHeader(200)
	_, _ = w.Write(resp)
}
