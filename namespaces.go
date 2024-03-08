package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func newNamespacesModel(
	size tea.WindowSizeMsg,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Namespace]{
		Title: renderTitle("select a namespace"),
		OnEnter: func(selected tui.Namespace, msgCh chan<- tea.Msg) {
			msgCh <- apisViewMsg{
				namespace: string(selected),
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
