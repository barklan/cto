package porter

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/barklan/cto/pkg/bot"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type PublicController struct {
	B *Base
}

func AllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func Serve(base *Base, s *bot.Sylon, queries chan<- QueryRequestWrap) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	ctrl := PublicController{B: base}

	oauthConf := initOAuth()

	r.Use(AllowCors)
	// TODO rewrite everything under publicController
	r.Route("/api/porter", func(r chi.Router) {
		r.Route("/query", func(r chi.Router) {
			r.Get("/exact", func(w http.ResponseWriter, r *http.Request) {
				serveLogExact(base, s, w, r)
			})
			r.Get("/range", func(w http.ResponseWriter, r *http.Request) {
				serveLogRange(base, s, queries, w, r)
			})
			r.Get("/poll", func(w http.ResponseWriter, r *http.Request) {
				pollLogRange(base, w, r)
			})
		})
		r.Route("/signin", func(r chi.Router) {
			r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
				handleOAuthLogin(base, oauthConf, w, r)
			})
			r.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
				handleOAuthCallback(base, oauthConf, w, r)
			})
		})
		// TODO chi does not get url path parameter otherwise
		r.Route("/project/{projectID}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				projectID := chi.URLParam(r, "projectID")
				ctrl.getProject(w, r, projectID)
			})
			r.Get("/issues", func(w http.ResponseWriter, r *http.Request) {
				projectID := chi.URLParam(r, "projectID")
				ctrl.projectIssues(w, r, projectID)
			})
			r.Get("/environments", func(w http.ResponseWriter, r *http.Request) {
				projectID := chi.URLParam(r, "projectID")
				ctrl.projectEnvs(w, r, projectID)
			})
			r.Get("/services", func(w http.ResponseWriter, r *http.Request) {
				projectID := chi.URLParam(r, "projectID")
				ctrl.projectServices(w, r, projectID)
			})
		})
		r.Route("/me/project", func(r chi.Router) {
			r.Use(render.SetContentType(render.ContentTypeJSON))
			r.Get("/", ctrl.getMyProjects)
			r.Get("/new", ctrl.newProjectRedirect)
		})
	})

	base.Log.Info("porter rest server is listening", zap.Int64("port", 9010))
	log.Panicln(http.ListenAndServe(":9010", r))
}
