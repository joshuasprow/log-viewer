package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/joshuasprow/log-viewer/k8s"
)

type Job struct {
	k8s.Job
}

func (j Job) FilterValue() string {
	return fmt.Sprintf("%s.%s", j.Namespace, j.Name)
}

func WrapJobs(jobs []k8s.Job) []list.Item {
	wrapped := make([]list.Item, len(jobs))
	for i, j := range jobs {
		wrapped[i] = Job{j}
	}
	return wrapped
}
