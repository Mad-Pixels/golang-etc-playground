package entrypoint

import (
	"net/http"
	"time"

	"github.com/Mad-Pixels/golang-playground/apps"
	"github.com/Mad-Pixels/golang-playground/apps/pkg/ws"
)

func customLogger(l ws.Logger) func(w http.ResponseWriter, r *http.Request) error {
	f := func(w http.ResponseWriter, r *http.Request) error {
		t := time.Now()
		l.Printf(
			"served request",
			map[string]any{
				"method":     r.Method,
				"path":       r.URL.String(),
				"replica_id": apps.ReplicaID(),
				"duration":   time.Since(t).String(),
				"request_id": r.Context().Value("uid"),
			},
		)
		return nil
	}
	return f
}
