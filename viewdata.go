package main

import (
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/tui"
)

type viewData struct {
	namespace        string
	api              tui.Api
	container        k8s.Container
	cronJob          k8s.CronJob
	cronJobJob       k8s.Job
	cronJobContainer k8s.Container
}
