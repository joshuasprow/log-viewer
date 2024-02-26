package pkg

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Pod struct {
	Namespace  string
	Name       string
	Containers []string
}

func GetPodsNext(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
) (
	[]Pod,
	error,
) {
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods := []Pod{}

	for _, item := range list.Items {
		containers := []string{}

		for _, container := range item.Spec.Containers {
			containers = append(containers, container.Name)
		}

		pods = append(pods, Pod{
			Namespace:  item.Namespace,
			Name:       item.Name,
			Containers: containers,
		})
	}

	return pods, nil
}
