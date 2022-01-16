package restcore

import (
	"encoding/json"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/barklan/cto/pkg/storage"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type restCore struct {
	Data *storage.Data
}

func Serve(data *storage.Data) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	ctrl := restCore{Data: data}

	r.Route("/api/core", func(r chi.Router) {
		r.Route("/debug/{project}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				getKey(data, w, r)
			})
		})

		r.Post("/setproject/{projectID}", func(w http.ResponseWriter, r *http.Request) {
			createNewProject(data, w, r)
		})

		r.Delete("/multi", ctrl.deletePrefix)
	})

	log.Panic(http.ListenAndServe(":8888", r))
}

func (c *restCore) deletePrefix(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	prefix := r.URL.Query().Get("prefix")

	var count uint64
	if err := c.Data.DB.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte(prefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := txn.Delete(k)
			if err != nil {
				return err
			}
			count++
		}
		return nil
	}); err != nil {
		c.Data.Log.Error("error when iterating badger keys", zap.Error(err))
		http.Error(w, "stopped prematurely", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := map[string]uint64{"count": count}
	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to marshal", http.StatusInternalServerError)
	}
	w.Write(respJson)
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
		data.Log.Error("failed to write debug response", zap.Error(err))
	}
}
