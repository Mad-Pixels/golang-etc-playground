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
	request.Name = r.Context().Value("uid").(string)

	///
	pod := k8s.Pod{
		Name:    request.Name,
		Image:   fmt.Sprintf("golang:%s-alpine3.18", request.Version),
		ExecCmd: []string{"go", "run", "/workspace/main.go"},
		Volumes: []*k8s.Volume{
			{
				Name: request.Name,
				Path: "/workspace",
				Files: []*k8s.File{
					{
						Filepath: "main.go",
						Source:   request.Source,
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
	client, err := k8s.SelfClient()
	if err != nil {
		panic(err)
	}
	for _, vol := range pod.Volumes {
		if _, err = vol.Create(r.Context(), client, "playground"); err != nil {
			panic(err)
		}
	}
	if _, err = pod.Create(r.Context(), client, "playground"); err != nil {
		panic(err)
	}

	watcher, err := pod.Watch(r.Context(), client, playgroundNs)
	if err != nil {
		panic(err)
	}
	for event := range watcher.ResultChan() {
		p, ok := event.Object.(*corev1.Pod)
		if !ok {
			fmt.Println("unexpected type")
			continue
		}
		fmt.Printf("Pod %s is in phase %s\n", p.Name, p.Status.Phase)
		if p.Status.Phase == "Failed" || p.Status.Phase == "Succeeded" {
			break
		}
	}
	output, err := pod.Logs(r.Context(), client, playgroundNs)
	if err != nil {
		panic(err)
	}
	responseOk(responseData{Data: output}, w, r)
}
