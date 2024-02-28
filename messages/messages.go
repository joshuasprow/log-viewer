package messages

import "github.com/joshuasprow/log-viewer/k8s"

type Error struct{ Err error }

func (e Error) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

type Namespaces []string

type Containers []k8s.Container

type CronJobs []string

type Logs []string
