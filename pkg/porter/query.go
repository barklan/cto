package porter

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/barklan/cto/pkg/bot"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type QStatus int

const (
	QWorking QStatus = iota
	QDone
	QFailed
)

type QResp struct {
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
		base.Log.Error("failed to set qmeta in cache", zap.Error(err))
	}
}

func checkProject(
	base *Base,
	w http.ResponseWriter,
	name,
	projectQ string,
) bool {
	var project string
	get := `--sql
		select project.id from project
		inner join client
			on project.client_id = client.id
		where client.email = $1 and project.id = $2`
	if err := base.R.Get(&project, get, name, projectQ); err != nil {
		http.Error(w, "Non existent project", 404)
		return false
	}
	if project != projectQ {
		return false
	}
	return true
}

func serveLogRange(
	base *Base,
	queries chan<- QueryRequestWrap,
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	rawQuery := r.URL.Query()

	tokenQ := rawQuery.Get("token")
	projectQ := rawQuery.Get("project")
	// FIXME check user here
	name, projectName, statusCode, ok := authorize(base, tokenQ)
	if !ok {
		w.WriteHeader(statusCode)
		return
	}
	if name != "guest" {
		if !checkProject(base, w, name, projectQ) {
			return
		}
		projectName = projectQ
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
		QID:       u4,
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
