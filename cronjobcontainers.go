package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func newCronJobContainersModel(
	size tea.WindowSizeMsg,
	cronJob k8s.CronJob,
	job k8s.Job,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Container]{
		Title: renderTitle(
			cronJob.Namespace,
			cronJob.Name,
			job.Name,
			"select a container",
		),
		OnEnter: func(selected tui.Container, msgCh chan<- tea.Msg) {
			msgCh <- cronJobLogsViewMsg{
				container: selected.Container,
			}
		},
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- cronJobJobsViewMsg{
				cronJob: cronJob,
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
