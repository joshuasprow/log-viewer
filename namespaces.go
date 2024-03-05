package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/messages"
	"github.com/joshuasprow/log-viewer/models"
)

type namespaceListItem string

func (n namespaceListItem) FilterValue() string {
	return string(n)
}

type namespacesModel struct {
	model *list.Model
	msgCh chan<- tea.Msg
}

func newNamespacesModel(
	size tea.WindowSizeMsg,
	msgCh chan<- tea.Msg,
) namespacesModel {
	m := models.DefaultListModel()
	m.Title = "namespaces"
	m.SetSize(size.Width, size.Height)

	return namespacesModel{
		model: &m,
		msgCh: msgCh,
	}
}

func (m namespacesModel) Init() tea.Cmd {
	return tea.Batch(
		m.model.StartSpinner(),
		func() tea.Msg {
			m.msgCh <- messages.Init{}
			return nil
		},
	)
}

func (m namespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- viewMsg{
				key:  apisKey,
				data: m.Selected(),
			}
		}
	case messages.Namespaces:
		items := []list.Item{}

		for _, namespace := range msg {
			items = append(items, namespaceListItem(namespace))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm

	return m, cmd
}

func (m namespacesModel) View() string {
	return m.model.View()
}

func (m namespacesModel) Selected() string {
	return m.model.SelectedItem().FilterValue()
}
