package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type logListItem string

func (i logListItem) FilterValue() string {
	return string(i)
}

type LogsModel struct {
	model     list.Model
	container ContainerListItem
}

func Logs(container ContainerListItem) *LogsModel {
	m := defaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "logs"

	return &LogsModel{
		model:     m,
		container: container,
	}
}

func (LogsModel) Init() tea.Cmd { return nil }

func (m *LogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrMsg:
		m.model.StopSpinner()
		return m, nil // todo: return error Cmd ?
	case LogsMsg:
		items := []list.Item{}

		for _, i := range msg {
			items = append(items, logListItem(i))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	var cmd tea.Cmd

	m.model, cmd = m.model.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

func (m *LogsModel) View() string {
	return m.model.View()
}

func (m *LogsModel) Selected() list.Item {
	return m.model.SelectedItem()
}
