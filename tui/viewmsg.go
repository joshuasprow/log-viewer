package tui

import (
	"github.com/joshuasprow/log-viewer/k8s"
)

type NamespacesViewMsg struct{}

type ApisViewMsg struct {
	Namespace string
}

type ContainersViewMsg struct {
	Namespace string
	Api       Api
}

type ContainerLogsViewMsg struct {
	Container k8s.Container
}

type CronJobsViewMsg struct {
	Namespace string
	Api       Api
}

type CronJobJobsViewMsg struct {
	CronJob k8s.CronJob
}

type CronJobContainersViewMsg struct {
	Job k8s.Job
}

type CronJobLogsViewMsg struct {
	Container k8s.Container
}
