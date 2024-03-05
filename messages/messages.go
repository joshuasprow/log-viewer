package messages

import "github.com/joshuasprow/log-viewer/k8s"

type Error struct{ Err error }

func (e Error) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

type Init struct{}

type Namespace struct {
	Name string
	Api  ApiKey
}

type Namespaces []string

type ApiKey string

const (
	ContainersApi ApiKey = "containers"
	CronJobsApi   ApiKey = "cronJobs"
)

func (k ApiKey) FilterValue() string {
	return string(k)
}

type Container k8s.Container
type Containers []k8s.Container

type CronJob k8s.CronJob
type CronJobs []k8s.CronJob

type Job k8s.Job
type Jobs []k8s.Job

type JobContainer k8s.Container
type JobContainers []k8s.Container

type Logs []string
