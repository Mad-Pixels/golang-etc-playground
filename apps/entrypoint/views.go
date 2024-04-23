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
	request.Name = r.Context().Value("uid").(string)

	mapSpecTpl, err := utils.Execute(
		playgroundMapSpec,
		request,
	)
	if err != nil {
		responseErrInternal(responseData{Message: err.Error()}, w, r)
		return
	}
	podSpecTpl, err := utils.Execute(
		playgroundPodSpec,
		request,
	)
	if err != nil {
		responseErrInternal(responseData{Message: err.Error()}, w, r)
		return
	}

	mapSpec, err := k8s.ToConfigMap(mapSpecTpl)
	if err != nil {
		responseErrInternal(responseData{Message: err.Error()}, w, r)
		return
	}
	podSpec, err := k8s.ToPod(podSpecTpl)
	if err != nil {
		responseErrInternal(responseData{Message: err.Error()}, w, r)
		return
	}

	_, err = k8s.ConfigMapCreate(r.Context(), playgroundNs, mapSpec)
	if err != nil {
		responseErrInternal(responseData{Message: err.Error()}, w, r)
		return
	}
	pod, err := k8s.PodCreate(r.Context(), playgroundNs, podSpec)
	if err != nil {
		responseErrInternal(responseData{Message: err.Error()}, w, r)
		return
	}
	responseOk(responseData{Data: pod.Name}, w, r)
}
