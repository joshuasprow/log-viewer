package k8s

import (
	"context"
	"fmt"
	"slices"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type Job struct {
	Namespace      string
	Name           string
	StartTime      time.Time
	CompletionTime time.Time
	Failed         int32
	Succeeded      int32
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
	list, err := clientset.
		BatchV1().
		Jobs(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}

	jobs := []Job{}

	for _, item := range list.Items {
		if slices.ContainsFunc(
			item.OwnerReferences,
			func(r metav1.OwnerReference) bool { return r.UID == cronJobUID },
		) {
			st := time.Time{}
			ct := time.Time{}

			if item.Status.StartTime != nil {
				st = item.Status.StartTime.Time
			}
			if item.Status.CompletionTime != nil {
				ct = item.Status.CompletionTime.Time
			}

			jobs = append(jobs, Job{
				Namespace:      item.Namespace,
				Name:           item.Name,
				StartTime:      st,
				CompletionTime: ct,
				Failed:         item.Status.Failed,
				Succeeded:      item.Status.Succeeded,
			})
		}
	}

	return jobs, nil
}
