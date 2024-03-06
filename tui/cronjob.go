package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/joshuasprow/log-viewer/k8s"
)

type CronJob struct {
	k8s.CronJob
}

func (c CronJob) FilterValue() string {
	return fmt.Sprintf("%s.%s", c.Namespace, c.Name)
}

func WrapCronJobs(cronJobs []k8s.CronJob) []list.Item {
	wrapped := make([]list.Item, len(cronJobs))
	for i, c := range cronJobs {
		wrapped[i] = CronJob{c}
	}
	return wrapped
}
