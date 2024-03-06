package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type CronJob struct {
	Namespace string
	UID       types.UID
	Name      string
	Jobs      []Job
}

func GetCronJobs(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
) (
	[]CronJob,
	error,
) {
	list, err := clientset.BatchV1().
		CronJobs(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	cronJobs := []CronJob{}

	for _, item := range list.Items {
		cronJobs = append(cronJobs, CronJob{
			Namespace: item.Namespace,
			UID:       item.UID,
			Name:      item.Name,
		})
	}

	return cronJobs, nil
}
