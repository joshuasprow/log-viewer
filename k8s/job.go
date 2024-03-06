package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type Job struct {
	Namespace string
	Name      string
}

func GetJobs(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	cronJobUID types.UID,
) (
	[]Job,
	error,
) {
	labelSelector := fmt.Sprintf("controller-uid=%s", string(cronJobUID))

	list, err := clientset.
		BatchV1().
		Jobs(namespace).
		List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}

	jobs := []Job{}

	for _, item := range list.Items {
		jobs = append(jobs, Job{
			Namespace: item.Namespace,
			Name:      item.Name,
		})
	}

	return jobs, nil
}
