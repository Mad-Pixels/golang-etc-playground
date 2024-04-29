package k8s

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	corev1 "k8s.io/api/core/v1"
)

func ToConfigMapSpec(raw string) (*corev1.ConfigMap, error) {
	var configMap corev1.ConfigMap
	if err := json.Unmarshal([]byte(raw), &configMap); err != nil {
		return nil, err
	}
	return &configMap, nil
}

func ToPodSpec(raw string) (*corev1.Pod, error) {
	var pod corev1.Pod
	if err := json.Unmarshal([]byte(raw), &pod); err != nil {
		return nil, err
	}
	return &pod, nil
}

func execTpl(data string, values any) (string, error) {
	funcMap := template.FuncMap{
		"escape": func(a string) string { return strings.ReplaceAll(a, `"`, `\"`) },
		"sub":    func(a, b int) int { return a - b },
		"toList": func(a interface{}) string {
			b, err := json.Marshal(a)
			if err != nil {
				return "[]"
			}
			return string(b)
		},
		"toJson": func(input string) string {
			b, err := json.Marshal(input)
			if err != nil {
				return ""
			}
			return string(b)
		},
	}
	tpl, err := template.New("tpl").Funcs(funcMap).Parse(data)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err = tpl.Execute(&body, values); err != nil {
		return "", err
	}
	return body.String(), nil
}
