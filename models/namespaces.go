package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func Namespaces(size tea.WindowSizeMsg, msgCh chan<- tea.Msg) tea.Model {
	options := defaults.ListModelOptions[tui.Namespace]{
		Title: tui.RenderTitle("select a namespace"),
		OnEnter: func(selected tui.Namespace, msgCh chan<- tea.Msg) {
			msgCh <- tui.ApisViewMsg{
				Namespace: string(selected),
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
