package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func Containers(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Container]{
		Title: tui.RenderTitle(namespace, "select a container"),
		OnEnter: func(selected tui.Container, msgCh chan<- tea.Msg) {
			msgCh <- tui.ContainerLogsViewMsg{
				Container: selected.Container,
			}
		},
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- tui.ApisViewMsg{
				Namespace: namespace,
			}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
