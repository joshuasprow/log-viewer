package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func CronJobLogs(
	size tea.WindowSizeMsg,
	cronJob k8s.CronJob,
	job k8s.Job,
	container k8s.Container,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Log]{
		Title: tui.RenderTitle(
			cronJob.Namespace,
			cronJob.Name,
			job.Name,
			container.Pod,
			container.Name,
			"logs",
		),
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- tui.CronJobContainersViewMsg{
				Job: job,
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
