package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/tui"
)

func newApisModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[tui.Api]{
		title: renderTitle(namespace, "select an API"),
		onEnter: func(selected tui.Api, msgCh chan<- tea.Msg) {
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
		onEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- namespacesViewMsg{}
		},
	}

	return newListModel(size, options, msgCh)
}
