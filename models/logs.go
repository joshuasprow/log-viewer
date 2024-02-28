package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/messages"
)

type logListItem string

func (i logListItem) FilterValue() string {
	return string(i)
}

type LogsModel struct {
	// a pointer is necessary for updating the spinner state
	model *list.Model
}

func Logs() LogsModel {
	m := defaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "logs"

	return LogsModel{model: &m}
}

func (LogsModel) Init() tea.Cmd { return nil }

func (m LogsModel) Update(msg tea.Msg) (LogsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.Logs:
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

func (m LogsModel) View() string {
	return m.model.View()
}

func (m LogsModel) Selected() list.Item {
	return m.model.SelectedItem()
}
