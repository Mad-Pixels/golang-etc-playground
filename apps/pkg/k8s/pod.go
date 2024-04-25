package k8s

import (
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return client.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{})
}

func PodWatch(ctx context.Context, namespace, name string) (watch.Interface, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}
	return client.CoreV1().Pods(namespace).Watch(ctx, metav1.SingleObject(metav1.ObjectMeta{
		Name: name,
	}))
}

// PodDelete delete a Pod in the specified namespace.
func PodDelete(ctx context.Context, namespace, name string) error {
	client, err := newClient()
	if err != nil {
		return err
	}
	return client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
