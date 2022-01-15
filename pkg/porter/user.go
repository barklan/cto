package porter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type ProjectResp struct {
	ID          string         `db:"id"`
	PrettyTitle sql.NullString `db:"pretty_title"`
	SecretKey   string         `db:"secret_key"`
}

func newProjectRedirect(base *Base, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	email, _, statusCode, ok := authorize(base, token)
	if !ok {
		w.WriteHeader(statusCode)
		return
	}
	pass := makeUserIntegrationPass(base, email)
	botName := "ctootestbot"
	if os.Getenv("CONFIG_ENV") == "prod" {
		botName = "ctoobot"
	}
	newProjectTGLink := fmt.Sprintf(
		"https://t.me/%s?startgroup=%s",
		botName,
		pass,
	)
	http.Redirect(w, r, newProjectTGLink, http.StatusTemporaryRedirect)
}

func getMyProjects(base *Base, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	name, _, statusCode, ok := authorize(base, token)
	if !ok {
		w.WriteHeader(statusCode)
		return
	}
	projects := make([]ProjectResp, 0)
	query := `--sql
		select project.id, pretty_title, secret_key from project
		inner join client
			on project.client_id = client.id
		where client.email = $1`
	if err := base.R.Select(&projects, query, name); err != nil {
		log.Println("err when selecting projects for user: ", err)
	}

	resp, err := json.Marshal(projects)
	if err != nil {
		http.Error(w, "Failed to marshal resp.", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}
