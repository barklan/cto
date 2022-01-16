package porter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/barklan/cto/pkg/storage/vars"
	"go.uber.org/zap"
)

type ProjectResp struct {
	ID          string         `db:"id"`
	SecretKey   string         `db:"secret_key"`
	PrettyTitle sql.NullString `db:"pretty_title"`
}

func (c *PublicController) newProjectRedirect(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	email, _, statusCode, ok := authorize(c.B, token)
	if !ok {
		w.WriteHeader(statusCode)
		return
	}
	pass := makeUserIntegrationPass(c.B, email)
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

func (c *PublicController) getMyProjects(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	name, _, statusCode, ok := authorize(c.B, token)
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
	if err := c.B.R.Select(&projects, query, name); err != nil {
		http.Error(w, "could not get project from db", 500)
		return
	}

	resp, err := json.Marshal(projects)
	if err != nil {
		http.Error(w, "Failed to marshal resp.", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func (c *PublicController) verifyProject(w http.ResponseWriter, email, projectID string) bool {
	var clientID string
	if err := c.B.R.Get(&clientID, `--sql
		select client_id from project
		where id = $1`,
		projectID,
	); err != nil {
		c.B.Log.Error("error when getting the owner of project", zap.String("project", projectID), zap.Error(err))
		http.Error(w, "no such project", http.StatusNotFound)
		return false
	}

	var ownerEmail string
	if err := c.B.R.Get(&ownerEmail, `--sql
		select email from client
		where id = $1`,
		clientID,
	); err != nil {
		http.Error(w, "could not find email for client", http.StatusInternalServerError)
		return false
	}
	if ownerEmail != email {
		http.Error(w, "you are not the owner", http.StatusForbidden)
	}
	return true
}

func (c *PublicController) projectStatus(w http.ResponseWriter, r *http.Request, projectID string) {
	token := r.URL.Query().Get("token")
	email, project, statusCode, ok := authorize(c.B, token)
	if !ok {
		w.WriteHeader(statusCode)
		return
	}
	if email == "guest" && project != projectID {
		http.Error(w, "project from token and path mismatch", http.StatusForbidden)
		return
	} else if email != "guest" && !c.verifyProject(w, email, projectID) {
		return
	}

	knownErrors, ok, err := c.B.Cache.GetVar(projectID, vars.KnownErrors)
	if err != nil {
		http.Error(w, "failed to get knownErorrs from cache", http.StatusInternalServerError)
	}
	if !ok {
		w.WriteHeader(404)
		_, _ = w.Write([]byte("no known errors"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(knownErrors)
}
