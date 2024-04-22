package entrypoint

import (
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/golang-playground/apps"
	"github.com/Mad-Pixels/golang-playground/apps/pkg/k8s"
	"github.com/Mad-Pixels/golang-playground/apps/pkg/utils"
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
	spec, err := utils.Execute(
		playgroundSpec,
		playgroundTmpl{
			Name:    r.Context().Value("uid").(string),
			Version: request.Version,
		},
	)
	if err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}

	podSpec, err := k8s.ToPod(spec)
	if err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}
	pod, err := k8s.PodCreate(r.Context(), playgroundNs, podSpec)
	if err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}
	responseOk(responseData{Data: pod.Name}, w, r)
}
