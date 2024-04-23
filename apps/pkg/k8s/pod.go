package k8s

import (
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ToPod parses a JSON string into a Pod object.
func ToPod(raw string) (*corev1.Pod, error) {
	var pod corev1.Pod
	if err := json.Unmarshal([]byte(raw), &pod); err != nil {
		return nil, err
	}
	return &pod, nil
}

// PodCreate creates a Pod in the specified namespace.
func PodCreate(ctx context.Context, namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
}
