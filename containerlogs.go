package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/tui"
)

func newContainerLogsModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[tui.Log]{
		onEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- containersViewMsg{
				namespace: namespace,
			}
		},
	}

	return newListModel(size, options, msgCh)
}
