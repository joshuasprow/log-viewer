package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models/defaults"
	"github.com/joshuasprow/log-viewer/tui"
)

func newApisModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := defaults.ListModelOptions[tui.Api]{
		Title: renderTitle(namespace, "select an API"),
		OnEnter: func(selected tui.Api, msgCh chan<- tea.Msg) {
			switch selected {
			case tui.ContainersApi:
				msgCh <- containersViewMsg{
					namespace: namespace,
					api:       selected,
				}
			case tui.CronJobsApi:
				msgCh <- cronJobsViewMsg{
					namespace: namespace,
					api:       selected,
				}
			}
		},
		OnEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- namespacesViewMsg{}
		},
	}

	return defaults.NewListModel(size, options, msgCh)
}
