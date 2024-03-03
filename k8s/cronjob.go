package k8s

import (
	"context"
	"fmt"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type Job struct {
	Namespace string
	Name      string
}

type CronJob struct {
	Namespace string
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
	cronJobsMap := map[types.UID]CronJob{}

	jlist, err := clientset.BatchV1().
		Jobs(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	for _, j := range jlist.Items {
		if len(j.OwnerReferences) == 0 {
			continue
		}

		job := Job{
			Namespace: j.Namespace,
			Name:      j.Name,
		}

		for _, owner := range j.OwnerReferences {
			if owner.Kind != "CronJob" {
				continue
			}

			cronJob, ok := cronJobsMap[owner.UID]
			if !ok {
				cronJob = CronJob{
					Namespace: j.Namespace,
					Name:      owner.Name,
					Jobs:      []Job{},
				}
			}

			cronJob.Jobs = append(cronJob.Jobs, job)

			cronJobsMap[owner.UID] = cronJob
		}
	}

	cronJobs := []CronJob{}

	for _, cj := range cronJobsMap {
		cronJobs = append(cronJobs, cj)
	}

	sort.Slice(cronJobs, func(i, j int) bool {
		return strings.Compare(cronJobs[i].Name, cronJobs[j].Name) < 0
	})

	return cronJobs, nil
}
