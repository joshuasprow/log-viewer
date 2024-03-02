package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CronJob struct {
	Namespace string
	Name      string
}

func GetCronJobs(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
) (
	[]string,
	error,
) {
	cronJobs, err := clientset.
		BatchV1().
		CronJobs(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list cronjobs: %w", err)
	}

	var names []string
	for _, cronJob := range cronJobs.Items {
		names = append(names, cronJob.Name)
	}

	return names, nil
}
