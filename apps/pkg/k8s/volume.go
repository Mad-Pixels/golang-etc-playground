package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var CfgMapSpecTpl = `
{
  "apiVersion": "v1",
  "kind": "ConfigMap",
  "metadata": {
    "name": "{{ .Name }}"
  },
  "data": { {{- $length := len .Files -}}
    {{- range $index, $file := .Files }}
    "{{ $file.Filepath }}": {{ toJson $file.Source }}{{if lt $index (sub $length 1)}},{{end}}
    {{- end }}
  }
}
`

type Volume struct {
	Name       string
	Path       string
	Files      []*File
	createOpts metav1.CreateOptions
	removeOpts metav1.DeleteOptions
}

type File struct {
	Filepath string
	Source   string
}

func (v *Volume) Spec() (*corev1.ConfigMap, error) {
	tpl, err := execTpl(CfgMapSpecTpl, v)
	if err != nil {
		return nil, err
	}
	return ToConfigMapSpec(tpl)
}

func (v *Volume) Create(ctx context.Context, client *kubernetes.Clientset, ns string) (*corev1.ConfigMap, error) {
	spec, err := v.Spec()
	if err != nil {
		return nil, err
	}
	return client.CoreV1().ConfigMaps(ns).Create(ctx, spec, v.createOpts)
}

func (v *Volume) Remove(ctx context.Context, client *kubernetes.Clientset, ns string) error {
	return client.CoreV1().ConfigMaps(ns).Delete(ctx, v.Name, v.removeOpts)
}
