package k8s

import (
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ToConfigMap parses a JSON string into a ConfigMap object.
func ToConfigMap(raw string) (*corev1.ConfigMap, error) {
	var configMap corev1.ConfigMap
	if err := json.Unmarshal([]byte(raw), &configMap); err != nil {
		return nil, err
	}
	return &configMap, nil
}

// ConfigMapCreate creates a ConfigMap in the specified namespace.
func ConfigMapCreate(ctx context.Context, namespace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, metav1.CreateOptions{})
}
