package porter

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/barklan/cto/pkg/bot"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Serve(base *Base, s *bot.Sylon, queries chan<- QueryRequestWrap) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api/porter", func(r chi.Router) {
		r.Route("/query", func(r chi.Router) {
			r.Get("/exact", func(w http.ResponseWriter, r *http.Request) {
				serveLogExact(base, s, w, r)
			})
			r.Post("/range", func(w http.ResponseWriter, r *http.Request) {
				serveLogRange(base, s, queries, w, r)
			})
			r.Get("/poll", func(w http.ResponseWriter, r *http.Request) {
				pollLogRange(base, w, r)
			})
		})
	})

	log.Panic(http.ListenAndServe(":9010", r))
}
