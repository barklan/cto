package porter

import (
	"fmt"
	"net/http"

	"github.com/barklan/cto/pkg/vars"
	"go.uber.org/zap"
)

func (c *PublicController) respFromCache(
	w http.ResponseWriter,
	r *http.Request,
	projectID, keyVar string,
) {
	resp, ok, err := c.B.Cache.GetVar(projectID, keyVar)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to get %s from cache", keyVar),
			http.StatusInternalServerError,
		)
		return
	}
	if !ok {
		http.Error(w, fmt.Sprintf("%s not found", keyVar), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		c.B.Log.Error(
			"failed to write response",
			zap.String("project", projectID),
			zap.String("keyVar", keyVar),
			zap.Error(err),
		)
	}
}

func (c *PublicController) projectIssues(w http.ResponseWriter, r *http.Request, projectID string) {
	if _, ok := c.guestOrOwner(w, r, projectID); !ok {
		return
	}
	c.respFromCache(w, r, projectID, vars.KnownErrors)
}

func (c *PublicController) projectEnvs(w http.ResponseWriter, r *http.Request, projectID string) {
	if _, ok := c.guestOrOwner(w, r, projectID); !ok {
		return
	}
	c.respFromCache(w, r, projectID, vars.KnownEnvs)
}

func (c *PublicController) projectServices(w http.ResponseWriter, r *http.Request, projectID string) {
	if _, ok := c.guestOrOwner(w, r, projectID); !ok {
		return
	}
	env := r.URL.Query().Get("env")
	c.respFromCache(w, r, projectID, env+vars.KnownServices)
}
