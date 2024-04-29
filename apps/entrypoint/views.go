package entrypoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mad-Pixels/golang-playground/apps"
	"github.com/Mad-Pixels/golang-playground/apps/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
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

	uid := r.Context().Value("uid").(string)
	source, err := request.SourceDecode()
	if err != nil {
		responseErrBadRequest(responseData{Message: "invalid body"}, w, r)
		return
	}
	client, err := k8s.SelfClient()
	if err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}
	pod := k8s.Pod{
		Name:    uid,
		Image:   fmt.Sprintf("golang:%s-alpine3.18", request.Version),
		ExecCmd: []string{"go", "run", "/workspace/main.go"},
		Volumes: []*k8s.Volume{
			{
				Name: uid,
				Path: "/workspace",
				Files: []*k8s.File{
					{
						Filepath: "main.go",
						Source:   source,
					},
				},
			},
		},
		Envs: []*k8s.Env{
			{
				Name:  "GOCACHE",
				Value: "/tmp/.cache",
			},
		},
	}
	for _, vol := range pod.Volumes {
		if _, err = vol.Create(r.Context(), client, "playground"); err != nil {
			fmt.Println(err)
			responseErrInternal(responseData{Message: "internal error"}, w, r)
			return
		}
	}
	if _, err = pod.Create(r.Context(), client, "playground"); err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}

	watcher, err := pod.Watch(r.Context(), client, playgroundNs)
	if err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}
	for event := range watcher.ResultChan() {
		p, ok := event.Object.(*corev1.Pod)
		if !ok {
			continue
		}
		if p.Status.Phase == "Failed" || p.Status.Phase == "Succeeded" {
			break
		}
	}
	output, err := pod.Logs(r.Context(), client, playgroundNs)
	if err != nil {
		responseErrInternal(responseData{Message: "internal error"}, w, r)
		return
	}
	defer func() {
		for _, vol := range pod.Volumes {
			_ = vol.Remove(r.Context(), client, playgroundNs)
		}
		_ = pod.Remove(r.Context(), client, playgroundNs)
	}()
	responseOk(responseData{Data: output}, w, r)
}
