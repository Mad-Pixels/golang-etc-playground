package entrypoint

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a App) routerTools() {
	router := func() chi.Router {
		r := chi.NewRouter()
		r.Get("/probe/liveness", handlerLivenessProbe)
		r.Get("/probe/readiness", handlerReadinessProbe)
		return r
	}
	a.server.Router().Mount("/api/internal", router())
}

func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func (a App) routerViews() {
	router := func() chi.Router {
		r := chi.NewRouter()
		r.Post("/playground", func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w, r)
			handlerPlayground(w, r)
		})
		r.Options("/playground", func(w http.ResponseWriter, r *http.Request) {
			setupCORS(&w, r)
			w.WriteHeader(http.StatusOK)
		})
		return r
	}
	a.server.Router().Mount("/api/v1", router())
}
