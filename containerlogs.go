package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
)

type logListItem string

func (i logListItem) FilterValue() string {
	return string(i)
}

type containerLogsModel struct {
	model     *list.Model
	msgCh     chan<- tea.Msg
	container k8s.Container
}

func newContainerLogsModel(
	size tea.WindowSizeMsg,
	container k8s.Container,
	msgCh chan<- tea.Msg,
) containerLogsModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = "container logs"

	return containerLogsModel{
		model:     &m,
		msgCh:     msgCh,
		container: container,
	}
}

func (m containerLogsModel) Init() tea.Cmd {
	return m.model.StartSpinner()
}

func (m containerLogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.msgCh <- containersViewMsg{namespace: m.container.Namespace}
			return m, nil
		}
	case containerLogsDataMsg:
		items := []list.Item{}

		for _, i := range msg {
			items = append(items, logListItem(i))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm

	return m, cmd
}

func (m containerLogsModel) View() string {
	return m.model.View()
}

func (m containerLogsModel) Selected() list.Item {
	return m.model.SelectedItem()
}
