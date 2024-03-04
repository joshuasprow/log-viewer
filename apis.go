package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/messages"
	"github.com/joshuasprow/log-viewer/models"
)

type apisListItem string

func (n apisListItem) FilterValue() string {
	return string(n)
}

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
	m.Title = "apis"
	m.SetItems([]list.Item{
		apisListItem(messages.ContainersApi),
		apisListItem(messages.CronJobsApi),
	})
	m.SetSize(size.Width, size.Height)

	return apisModel{
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
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
		case "enter":
			m.msgCh <- messages.Namespace{
				Name: m.namespace,
				Api:  messages.Api(m.Selected().FilterValue()),
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

func (m apisModel) Selected() list.Item {
	return m.model.SelectedItem()
}
