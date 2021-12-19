package porter

import (
	"log"
	"net/http"

	"github.com/barklan/cto/pkg/bot"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Serve(base *Base, s *bot.Sylon, reqs chan<- QueryRequest) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api/porter", func(r chi.Router) {
		r.Route("/query/exact", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				serveLogExact(base, s, w, r)
			})
		})
	})

	log.Panic(http.ListenAndServe(":9010", r))
}
