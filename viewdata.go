package main

import (
	"github.com/joshuasprow/log-viewer/k8s"
)

type viewData struct {
	namespace        string
	container        k8s.Container
	cronJob          k8s.CronJob
	cronJobJob       k8s.Job
	cronJobContainer k8s.Container
}
