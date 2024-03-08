package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func newCronJobsModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.CronJob]{
		ShowDescription: true,
		Title:           renderTitle(namespace, "select a cron job"),
		OnEnter: func(selected tui.CronJob, msgCh chan<- tea.Msg) {
			msgCh <- cronJobJobsViewMsg{
				cronJob: selected.CronJob,
			}
		},
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- apisViewMsg{
				namespace: namespace,
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
