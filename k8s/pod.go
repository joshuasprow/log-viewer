package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPods(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
) (
	[]string,
	error,
) {
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods := []string{}

	for _, item := range list.Items {
		pods = append(pods, item.Name)
	}

	return pods, nil
}