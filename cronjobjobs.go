package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/tui"
)

func newCronJobJobsModel(
	size tea.WindowSizeMsg,
	cronJob k8s.CronJob,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[tui.Job]{
		title: renderTitle(cronJob.Namespace, cronJob.Name, "select a job"),
		onEnter: func(selected tui.Job, msgCh chan<- tea.Msg) {
			msgCh <- cronJobContainersViewMsg{
				job: selected.Job,
			}
		},
		onEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- cronJobsViewMsg{
				namespace: cronJob.Namespace,
			}
		},
	}

	return newListModel(size, options, msgCh)
}
