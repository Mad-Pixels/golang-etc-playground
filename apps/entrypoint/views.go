package entrypoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mad-Pixels/golang-playground/apps"
)

func handlerLivenessProbe(w http.ResponseWriter, r *http.Request) {
	responseOk(responseData{Host: apps.ReplicaID()}, w, r)
}

func handlerReadinessProbe(w http.ResponseWriter, r *http.Request) {
	responseOk(responseData{Host: apps.ReplicaID()}, w, r)
}

func handlerPlayground(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseErrBadRequest(responseData{Message: "only POST method allowed"}, w, r)
		return
	}
	var request requestPlayground
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responseErrBadRequest(responseData{Message: "invalid body"}, w, r)
		return
	}

	//k8s.Pod()
	responseOk(responseData{Data: fmt.Sprintf("v%s-%s", request.Version, r.Context().Value("uid"))}, w, r)
}
