package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/joshuasprow/log-viewer/k8s"
)

type CronJob struct {
	k8s.CronJob
}

func (c CronJob) Title() string {
	return fmt.Sprintf("%s.%s", c.Namespace, c.Name)
}

func (c CronJob) Description() string {
	return fmt.Sprintf(
		"last_scheduled=%s",
		c.LastScheduleTime.Format("2006-01-02T15:04:05"),
	)
}

func (c CronJob) FilterValue() string {
	return c.Title()
}

func WrapCronJobs(cronJobs []k8s.CronJob) []list.Item {
	wrapped := make([]list.Item, len(cronJobs))
	for i, c := range cronJobs {
		wrapped[i] = CronJob{c}
	}
	return wrapped
}
