package k8s

import (
	"context"
	"strings"

	"github.com/joshuasprow/log-viewer/pkg"
	v1 "k8s.io/api/core/v1"
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

func GetPodLogs(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	pod string,
	container string,
) (
	[]string,
	error,
) {
	data, err := clientset.
		CoreV1().
		Pods(namespace).
		GetLogs(
			pod,
			&v1.PodLogOptions{
				Container: container,
				TailLines: pkg.Ptr[int64](20),
			},
		).
		Do(ctx).
		Raw()
	if err != nil {
		return nil, err
	}

	logs := []string{}

	for _, line := range strings.Split(string(data), "\n") {
		l := strings.TrimSpace(line)

		if l == "" {
			continue
		}

		logs = append(logs, l)
	}

	return logs, nil
}
