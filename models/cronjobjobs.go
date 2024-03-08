package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func CronJobJobs(
	size tea.WindowSizeMsg,
	cronJob k8s.CronJob,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Job]{
		ShowDescription: true,
		Title: tui.RenderTitle(
			cronJob.Namespace,
			cronJob.Name,
			"select a job",
		),
		OnEnter: func(selected tui.Job, msgCh chan<- tea.Msg) {
			msgCh <- tui.CronJobContainersViewMsg{
				Job: selected.Job,
			}
		},
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- tui.CronJobsViewMsg{
				Namespace: cronJob.Namespace,
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
