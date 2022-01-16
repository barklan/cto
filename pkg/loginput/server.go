package loginput

import (
	"crypto/subtle"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/barklan/cto/pkg/postgres/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

func authorizeRequest(rdb *sqlx.DB, r *http.Request) (string, bool) {
	projectName, password, ok := r.BasicAuth()
	if !ok {
		log.Println("error parsing basic auth")
		return "", false
	}

	project := models.Project{}
	if err := rdb.Get(&project, "select * from project where id = $1", projectName); err != nil {
		log.Println("project not found from basic auth")
		return "", false
	}

	if subtle.ConstantTimeCompare([]byte(password), []byte(project.SecretKey)) == 1 {
		return project.ID, true
	}

	return "", false
}

func logOneRequest(
	lg *zap.Logger,
	rdb *sqlx.DB,
	reqs chan<- LogRequest,
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	projectName, ok := authorizeRequest(rdb, r)
	if !ok {
		log.Warn("recieved unauthorized request")
		w.WriteHeader(401)
		return
	}

	lg.Info("recieved log dump", zap.String("project", projectName))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		lg.Warn("failed to read loginput body of project", zap.String("project", projectName))
		w.WriteHeader(400)
		return
	}

	reqs <- LogRequest{
		ProjectID: projectName,
		Body:      body,
	}
}

func Serve(lg *zap.Logger, rdb *sqlx.DB, reqs chan<- LogRequest) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/loginput/fluentd", func(w http.ResponseWriter, r *http.Request) {
		logOneRequest(lg, rdb, reqs, w, r)
	})
	log.Panic(http.ListenAndServe(":8900", r))
}
