package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	m.Title = renderTitle("select a namespaces")
	m.SetSize(size.Width, size.Height)

	return namespacesModel{
		model: &m,
		msgCh: msgCh,
	}
}

func newNamespacesModelNext(
	size tea.WindowSizeMsg,
	msgCh chan<- tea.Msg,
) tea.Model {
	options := listModelOptions[namespaceListItem]{
		onEnter: func(selected namespaceListItem, msgCh chan<- tea.Msg) {
			msgCh <- apisViewMsg{
				namespace: string(selected),
			}
		},
	}

	return newListModel(size, options, msgCh)
}

func (m namespacesModel) Init() tea.Cmd {
	return m.model.StartSpinner()
}

func (m namespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- apisViewMsg{namespace: m.Selected()}
			return m, nil
		}
	case namespacesDataMsg:
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
