package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func Pod(ctx context.Context) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// Создание клиента для взаимодействия с Kubernetes API.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Описание pod.
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "demo-container",
					Image: "nginx",
				},
			},
		},
	}

	// Создание pod в default namespace.
	_, err = clientset.CoreV1().Pods("playground").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Pod created successfully!")
}
