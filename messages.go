package main

import "github.com/joshuasprow/log-viewer/k8s"

type viewMsg struct {
	key  viewKey
	data any
}

type namespacesDataMsg []string
type containersDataMsg []k8s.Container
type containerLogsDataMsg []string
type cronJobsDataMsg []k8s.CronJob
type cronJobContainersDataMsg []k8s.Container
type cronJobLogsDataMsg []string
