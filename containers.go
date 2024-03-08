package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func newContainersModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Container]{
		Title: renderTitle(namespace, "select a container"),
		OnEnter: func(selected tui.Container, msgCh chan<- tea.Msg) {
			msgCh <- containerLogsViewMsg{
				container: selected.Container,
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
