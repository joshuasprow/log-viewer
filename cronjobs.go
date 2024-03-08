package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/tui"
)

func newCronJobsModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[tui.CronJob]{
		showDescription: true,
		title:           renderTitle(namespace, "select a cron job"),
		onEnter: func(selected tui.CronJob, msgCh chan<- tea.Msg) {
			msgCh <- cronJobJobsViewMsg{
				cronJob: selected.CronJob,
			}
		},
		onEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- apisViewMsg{
				namespace: namespace,
			}
		},
	}

	return newListModel(size, options, msgCh)
}
