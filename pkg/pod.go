package pkg

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPods(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
) (
	[]v1.Pod,
	error,
) {
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods := []v1.Pod{}

	for _, item := range list.Items {
		pods = append(pods, item)
	}

	return pods, nil
}
