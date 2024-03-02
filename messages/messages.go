package messages

import "github.com/joshuasprow/log-viewer/k8s"

type Error struct{ Err error }

func (e Error) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

type Api string

const (
	ContainersApi Api = "containers"
	CronJobsApi   Api = "cronJobs"
)

type Namespace struct {
	Name string
	Api  Api
}

type Namespaces []string

type Container k8s.Container
type Containers []k8s.Container

type CronJob k8s.CronJob
type CronJobs []k8s.CronJob

type Logs []string
