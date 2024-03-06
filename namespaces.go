package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/tui"
)

func newNamespacesModel(
	size tea.WindowSizeMsg,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[tui.Namespace]{
		onEnter: func(selected tui.Namespace, msgCh chan<- tea.Msg) {
			msgCh <- apisViewMsg{
				namespace: string(selected),
			}
		},
	}

	return newListModel(size, options, msgCh)
}
