package k8s

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

type Container struct {
	Namespace string
	Pod       string
	Name      string
}

func GetContainers(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	labelSelector string,
) (
	[]Container,
	error,
) {
	pods, err := GetPods(ctx, clientset, namespace, labelSelector)
	if err != nil {
		return nil, fmt.Errorf("load model data: %w", err)
	}

	containers := []Container{}

	for _, pod := range pods {
		if len(pod.Spec.Containers) == 0 {
			containers = append(containers, Container{
				Namespace: pod.Namespace,
				Pod:       pod.Name,
			})
			continue
		}

		for _, container := range pod.Spec.Containers {
			containers = append(containers, Container{
				Namespace: pod.Namespace,
				Pod:       pod.Name,
				Name:      container.Name,
			})
		}
	}

	return containers, nil
}
