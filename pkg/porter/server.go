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

	oauthConf := initOAuth()

	r.Route("/api/porter", func(r chi.Router) {
		r.Route("/query", func(r chi.Router) {
			r.Get("/exact", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				serveLogExact(base, s, w, r)
			})
			r.Get("/range", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				serveLogRange(base, s, queries, w, r)
			})
			r.Get("/poll", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				pollLogRange(base, w, r)
			})
		})
		r.Route("/signin", func(r chi.Router) {
			r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				handleOAuthLogin(base, oauthConf, w, r)
			})
			r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				handleOAuthCallback(base, oauthConf, w, r)
			})
		})
		r.Route("/me/project", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				getMyProjects(base, w, r)
			})
			r.Get("/new", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				newProjectRedirect(base, w, r)
			})
		})
	})

	log.Println("porter rest server listening on 9010")
	log.Panic(http.ListenAndServe(":9010", r))
}
