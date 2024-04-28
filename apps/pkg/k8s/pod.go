package k8s

import (
	"bytes"
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

var PodSpecTpl = `
{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "name": "{{ .Name }}"
    },
    "spec": {
        "restartPolicy": "Never",
        "containers": [
            {
                "name": "{{ .Name }}",
                "image": "{{ .Image }}",
                "securityContext": {
                    "runAsNonRoot": true,
                    "runAsUser": 1000
                },
                "env": [{{- $len := len .Envs }}{{ range $index, $env := .Envs }}
                    {
                        "name": "{{ $env.Name }}",
                        "value": "{{ $env.Value }}"{{ if lt $index (sub $len 1) }},
                        {{- end }}
                    }{{ end }}
                ],
                "resources": {
                    "requests": {
                        "cpu": "1000m",
                        "memory": "400Mi"
                    },
                    "limits": {
                        "cpu": "2000m",
                        "memory": "800Mi"
                    }
                },
                "command": {{ toJson .ExecCmd }},
                "volumeMounts": [
                    {{- $len := len .Volumes }}{{ range $index, $vol := .Volumes }}
                    {
                        "name": "{{ $vol.Name }}",
                        "mountPath": "{{ $vol.Path }}"{{ if lt $index (sub $len 1) }},
                        {{- end }}
                    }{{ end }}
                ]
            }
        ],
        "volumes": [
            {{- $len := len .Volumes }}{{ range $index, $vol := .Volumes }}
            {
                "name": "{{ $vol.Name }}",
                "configMap": {
                    "name": "{{ $vol.Name }}"
                }{{ if lt $index (sub $len 1) }},
                {{- end }}
            }{{ end }}
        ]
    }
}
`

type Env struct {
	Name  string
	Value string
}

type Pod struct {
	Name    string
	Image   string
	ExecCmd []string
	Volumes []*Volume
	Envs    []*Env

	createOpts metav1.CreateOptions
	removeOpts metav1.DeleteOptions
	logOpts    v1.PodLogOptions
}

func (p *Pod) Spec() (*corev1.Pod, error) {
	tpl, err := execTpl(PodSpecTpl, p)
	if err != nil {
		return nil, err
	}
	return ToPodSpec(tpl)
}

func (p *Pod) Create(ctx context.Context, client *kubernetes.Clientset, ns string) (*corev1.Pod, error) {
	spec, err := p.Spec()
	if err != nil {
		return nil, err
	}
	return client.CoreV1().Pods(ns).Create(ctx, spec, p.createOpts)
}

func (p *Pod) Remove(ctx context.Context, client *kubernetes.Clientset, ns string) error {
	return client.CoreV1().Pods(ns).Delete(ctx, p.Name, p.removeOpts)
}

func (p *Pod) Watch(ctx context.Context, client *kubernetes.Clientset, ns string) (watch.Interface, error) {
	return client.CoreV1().Pods(ns).Watch(
		ctx,
		metav1.SingleObject(metav1.ObjectMeta{Name: p.Name}),
	)
}

func (p *Pod) Logs(ctx context.Context, client *kubernetes.Clientset, ns string) (string, error) {
	req := client.CoreV1().Pods(ns).GetLogs(p.Name, &p.logOpts)
	logs, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}

	defer logs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
