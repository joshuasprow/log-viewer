package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/tui"
)

func newCronJobLogsModel(
	size tea.WindowSizeMsg,
	cronJob k8s.CronJob,
	job k8s.Job,
	container k8s.Container,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[tui.Log]{
		title: renderTitle(
			cronJob.Namespace,
			cronJob.Name,
			job.Name,
			container.Pod,
			container.Name,
			"logs",
		),
		onEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- cronJobContainersViewMsg{
				job: job,
			}
		},
	}

	return newListModel(size, options, msgCh)
}
