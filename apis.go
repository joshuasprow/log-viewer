package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models"
)

type apisModel struct {
	model     *list.Model
	namespace string
	msgCh     chan<- tea.Msg
}

func newApisViewModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) apisModel {
	m := models.DefaultListModel()
	m.SetItems([]list.Item{
		containersKey,
		cronJobsKey,
	})
	m.SetSize(size.Width, size.Height)
	m.Title = renderTitle(namespace, "select an api")

	return apisModel{
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
}

func newApisModelNext(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[viewKey]{
		onEnter: func(selected viewKey, msgCh chan<- tea.Msg) {
			switch selected {
			case containersKey:
				msgCh <- containersViewMsg{namespace: namespace}
			case cronJobsKey:
				msgCh <- cronJobsViewMsg{namespace: namespace}
			}
		},
		onEsc: func(msgCh chan<- tea.Msg) {
			msgCh <- namespacesViewMsg{}
		},
	}

	return newListModel(size, options, msgCh)
}

func (m apisModel) Init() tea.Cmd {
	return nil
}

func (m apisModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.msgCh <- namespacesViewMsg{}
			return m, nil
		case "enter":
			switch m.Selected() {
			case containersKey:
				m.msgCh <- containersViewMsg{namespace: m.namespace}
				return m, nil
			case cronJobsKey:
				m.msgCh <- cronJobsViewMsg{namespace: m.namespace}
				return m, nil
			}
		}
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm

	return m, cmd
}

func (m apisModel) View() string {
	return m.model.View()
}

func (m apisModel) Selected() viewKey {
	return m.model.SelectedItem().(viewKey)
}
