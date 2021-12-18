package restcore

import (
	"log"
	"net/http"

	"github.com/barklan/cto/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func createNewProject(data *storage.Data, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	projectID := chi.URLParam(r, "projectID")
	data.SetVar(projectID, "", "", -1)
	w.WriteHeader(http.StatusOK)
}

func Serve(data *storage.Data) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/api/core/setproject/{projectID}", func(w http.ResponseWriter, r *http.Request) {
		createNewProject(data, w, r)
	})
	log.Panic(http.ListenAndServe(":8888", r))
}
