package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/tui"
)

func newContainersModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[tui.Container]{
		title: renderTitle(namespace, "select a container"),
		onEnter: func(selected tui.Container, msgCh chan<- tea.Msg) {
			msgCh <- containerLogsViewMsg{
				container: selected.Container,
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
