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
	[]CronJob,
	error,
) {
	cronJobs, err := clientset.
		BatchV1().
		CronJobs(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list cronjobs: %w", err)
	}

	jobs := []CronJob{}

	for _, cronJob := range cronJobs.Items {
		jobs = append(jobs, CronJob{
			Namespace: cronJob.Namespace,
			Name:      cronJob.Name,
		})
	}

	return jobs, nil
}
