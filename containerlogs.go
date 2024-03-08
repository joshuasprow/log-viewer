package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func newContainerLogsModel(
	size tea.WindowSizeMsg,
	container k8s.Container,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Log]{
		Title: renderTitle(
			container.Namespace,
			container.Pod,
			container.Name,
			"logs",
		),
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- containersViewMsg{
				namespace: container.Namespace,
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
