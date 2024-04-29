package entrypoint

import (
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/golang-etc-playground/apps"
)

type responseData struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Host    string `json:"host"`
	Id      string `json:"id"`
	Status  int    `json:"status"`
}

func responseOk(data responseData, w http.ResponseWriter, r *http.Request) {
	data.Status = http.StatusOK
	response(data, w, r)
}

func responseErrInternal(data responseData, w http.ResponseWriter, r *http.Request) {
	data.Status = http.StatusInternalServerError
	response(data, w, r)
}

func responseErrBadRequest(data responseData, w http.ResponseWriter, r *http.Request) {
	data.Status = http.StatusBadRequest
	response(data, w, r)
}

func response(data responseData, w http.ResponseWriter, r *http.Request) {
	data.Id = r.Context().Value("uid").(string)
	data.Host = apps.ReplicaID()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(data.Status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		data.Data = nil
		responseErrInternal(data, w, r)
	}
}
