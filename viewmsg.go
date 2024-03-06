package main

import "github.com/joshuasprow/log-viewer/k8s"

type namespacesViewMsg struct{}
type apisViewMsg struct{ namespace string }
type containersViewMsg struct{ namespace string }
type containerLogsViewMsg struct{ container k8s.Container }
type cronJobsViewMsg struct{ namespace string }
type cronJobJobsViewMsg struct{ cronJob k8s.CronJob }
type cronJobContainersViewMsg struct{ job k8s.Job }
type cronJobLogsViewMsg struct{ container k8s.Container }

type namespacesDataMsg []string
type containersDataMsg []k8s.Container
type containerLogsDataMsg []string
type cronJobsDataMsg []k8s.CronJob
type cronJobContainersDataMsg []k8s.Container
type cronJobLogsDataMsg []string
