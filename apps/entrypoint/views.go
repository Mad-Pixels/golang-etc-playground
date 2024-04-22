package entrypoint

import (
	"encoding/json"
	"net/http"

	apps "github.com/Mad-Pixels/golang-playground"
)

type probeResponse struct {
	Host    string `json:"host"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status"`
}

func handlerInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome"))
}

func handlerLivenessProbe(w http.ResponseWriter, _ *http.Request) {
	response := probeResponse{
		Status:  "OK",
		Message: "LivenessProbe",
		Host:    apps.ReplicaID(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func handlerReadnessProbe(w http.ResponseWriter, _ *http.Request) {
	response := probeResponse{
		Status:  "OK",
		Message: "ReadnessProbe",
		Host:    apps.ReplicaID(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
