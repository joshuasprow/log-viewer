package main

import (
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/tui"
)

type namespacesViewMsg struct{}
type apisViewMsg struct{ namespace string }
type containersViewMsg struct {
	namespace string
	api       tui.Api
}
type containerLogsViewMsg struct{ container k8s.Container }
type cronJobsViewMsg struct {
	namespace string
	api       tui.Api
}
type cronJobJobsViewMsg struct{ cronJob k8s.CronJob }
type cronJobContainersViewMsg struct{ job k8s.Job }
type cronJobLogsViewMsg struct{ container k8s.Container }
