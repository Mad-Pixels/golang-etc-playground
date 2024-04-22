package entrypoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Mad-Pixels/golang-playground/apps"
	"github.com/Mad-Pixels/golang-playground/apps/pkg/k8s"
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

	tmpl, err := template.New("pod").Parse(playgroundSpec)
	if err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}
	var tmplSpec bytes.Buffer
	err = tmpl.Execute(
		&tmplSpec,
		playgroundTmpl{
			Name:    r.Context().Value("uid").(string),
			Version: request.Version,
		},
	)
	if err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}

	spec, err := k8s.ToPod(tmplSpec.String())
	if err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}
	if _, err := k8s.PodCreate(r.Context(), playgroundNs, spec); err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}
	responseOk(responseData{Data: fmt.Sprintf("v%s-%s", request.Version, r.Context().Value("uid"))}, w, r)
}
