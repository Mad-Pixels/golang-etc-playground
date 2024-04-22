package entrypoint

import (
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

func (a App) routerViews() {
	router := func() chi.Router {
		r := chi.NewRouter()
		r.Get("/", handlerInfo)
		return r
	}
	a.server.Router().Mount("/api/v1", router())
}
