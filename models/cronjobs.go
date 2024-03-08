package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func CronJobs(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.CronJob]{
		ShowDescription: true,
		Title: tui.RenderTitle(
			namespace,
			"select a cron job",
		),
		OnEnter: func(selected tui.CronJob, msgCh chan<- tea.Msg) {
			msgCh <- tui.CronJobJobsViewMsg{
				CronJob: selected.CronJob,
			}
		},
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- tui.ApisViewMsg{
				Namespace: namespace,
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
