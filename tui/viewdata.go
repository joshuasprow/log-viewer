package tui

import (
	"github.com/joshuasprow/log-viewer/k8s"
)

type ViewData struct {
	Namespace        string
	Api              Api
	Container        k8s.Container
	CronJob          k8s.CronJob
	CronJobJob       k8s.Job
	CronJobContainer k8s.Container
}
