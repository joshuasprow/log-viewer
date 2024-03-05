package main

import (
	"fmt"

	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
)

type viewKey string

const (
	namespacesKey viewKey = "namespaces"
	apisKey       viewKey = namespacesKey + ".apis"

	containersKey    viewKey = apisKey + ".containers"
	containerLogsKey viewKey = containersKey + ".logs"

	cronJobsKey    viewKey = apisKey + ".cronjobs"
	cronJobJobsKey viewKey = cronJobsKey + ".jobs"
	cronJobLogsKey viewKey = cronJobJobsKey + ".logs"
)

func (k viewKey) FilterValue() string {
	return string(k)
}

type viewData struct {
	namespace        string
	api              messages.ApiKey
	container        k8s.Container
	cronJob          k8s.CronJob
	cronJobContainer k8s.Container
}

func updateViewData(data viewData, msg viewMsg) (viewData, error) {
	switch msg.key {
	case namespacesKey:
		data.namespace = msg.data.(string)
	case apisKey:
		data.namespace = msg.data.(string)
	case containersKey:
		data.namespace = msg.data.(string)
	case containerLogsKey:
		data.container = msg.data.(k8s.Container)
	case cronJobsKey:
		data.namespace = msg.data.(string)
	case cronJobJobsKey:
		data.cronJob = msg.data.(k8s.CronJob)
	case cronJobLogsKey:
		data.cronJobContainer = msg.data.(k8s.Container)
	default:
		return data, fmt.Errorf("unknown view key: %s", msg.key)
	}

	return data, nil
}

type viewMsg struct {
	key  viewKey
	data any
}
