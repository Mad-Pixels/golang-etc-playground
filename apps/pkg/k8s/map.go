package k8s

import (
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return client.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, metav1.CreateOptions{})
}

// ConfigMapDelete deletes a ConfigMap in the specified namespace.
func ConfigMapDelete(ctx context.Context, namespace, name string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	deleteOptions := metav1.DeleteOptions{}
	return client.CoreV1().ConfigMaps(namespace).Delete(ctx, name, deleteOptions)
}
