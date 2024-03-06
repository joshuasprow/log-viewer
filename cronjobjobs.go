package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/tui"
)

func newCronJobJobsModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[tui.Job]{
		onEnter: func(selected tui.Job, msgCh chan<- tea.Msg) {
			msgCh <- cronJobContainersViewMsg{
				job: selected.Job,
			}
		},
		onEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- cronJobsViewMsg{
				namespace: namespace,
			}
		},
	}

	return newListModel(size, options, msgCh)
}
