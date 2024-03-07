package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/joshuasprow/log-viewer/k8s"
)

type Job struct {
	k8s.Job
}

func (j Job) Title() string {
	icon := "✅"
	if j.Failed > 0 {
		icon = "🚫"
	}
	return fmt.Sprintf("%s %s.%s", icon, j.Namespace, j.Name)
}

func (j Job) Description() string {
	return fmt.Sprintf(
		"start_time=%s failed=%d succeeded=%d",
		j.StartTime.Format("2006-01-02T15:04:05"),
		j.Failed,
		j.Succeeded,
	)
}

func (j Job) FilterValue() string {
	return j.Title()
}

func WrapJobs(jobs []k8s.Job) []list.Item {
	wrapped := make([]list.Item, len(jobs))
	for i, j := range jobs {
		wrapped[i] = Job{j}
	}
	return wrapped
}
