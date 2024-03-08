package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func Apis(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Api]{
		Title: tui.RenderTitle(namespace, "select an API"),
		OnEnter: func(selected tui.Api, msgCh chan<- tea.Msg) {
			switch selected {
			case tui.ContainersApi:
				msgCh <- tui.ContainersViewMsg{
					Namespace: namespace,
					Api:       selected,
				}
			case tui.CronJobsApi:
				msgCh <- tui.CronJobsViewMsg{
					Namespace: namespace,
					Api:       selected,
				}
			}
		},
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- tui.NamespacesViewMsg{}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
