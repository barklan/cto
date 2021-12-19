package restcore

import (
	"log"
	"net/http"

	"github.com/barklan/cto/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Serve(data *storage.Data) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api/core", func(r chi.Router) {
		r.Route("/debug/{project}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				getKey(data, w, r)
			})
		})

		r.Post("/setproject/{projectID}", func(w http.ResponseWriter, r *http.Request) {
			createNewProject(data, w, r)
		})
	})

	log.Panic(http.ListenAndServe(":8888", r))
}

func createNewProject(data *storage.Data, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	projectID := chi.URLParam(r, "projectID")
	data.SetVar(projectID, "", "!", -1)
	w.WriteHeader(http.StatusOK)
}

func getKey(data *storage.Data, w http.ResponseWriter, r *http.Request) {
	project := chi.URLParam(r, "project")
	key := r.URL.Query().Get("key")

	if !data.VarExists(project, key) {
		http.Error(w, "badger variable not found", 404)
		return
	}

	val := data.GetVar(project, key)

	_, err := w.Write(val)
	if err != nil {
		log.Println("failed to write debug response")
	}
}
