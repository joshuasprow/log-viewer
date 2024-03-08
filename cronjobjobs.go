package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func newCronJobJobsModel(
	size tea.WindowSizeMsg,
	cronJob k8s.CronJob,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Job]{
		ShowDescription: true,
		Title: renderTitle(
			cronJob.Namespace,
			cronJob.Name,
			"select a job",
		),
		OnEnter: func(selected tui.Job, msgCh chan<- tea.Msg) {
			msgCh <- cronJobContainersViewMsg{
				job: selected.Job,
			}
		},
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- cronJobsViewMsg{
				namespace: cronJob.Namespace,
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
