package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/tui"
)

func newCronJobLogsModel(
	size tea.WindowSizeMsg,
	job k8s.Job,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[tui.Log]{
		onEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- cronJobContainersViewMsg{
				job: job,
			}
		},
	}

	return newListModel(size, options, msgCh)
}
