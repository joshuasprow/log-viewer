package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/models"
)

type listModel[ItemType any] struct {
	model   *list.Model
	options listModelOptions[ItemType]
	msgCh   chan<- tea.Msg
}

type listModelOptions[ItemType any] struct {
	onEnter         func(selected ItemType, msgCh chan<- tea.Msg)
	onEsc           func(msgCh chan<- tea.Msg)
	showDescription bool
	title           string
}

func newListModel[ItemType any](
	size tea.WindowSizeMsg,
	options listModelOptions[ItemType],
	msgCh chan<- tea.Msg,
) listModel[ItemType] {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = options.title

	return listModel[ItemType]{
		model:   &m,
		options: options,
		msgCh:   msgCh,
	}
}

func (m listModel[ItemType]) Init() tea.Cmd {
	return m.model.StartSpinner()
}

func (m listModel[ItemType]) Update(
	msg tea.Msg,
) (
	tea.Model,
	tea.Cmd,
) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			if m.options.onEsc != nil {
				m.options.onEsc(m.msgCh)
				return m, nil
			}
		case "enter":
			if m.options.onEnter != nil {
				m.options.onEnter(m.Selected(), m.msgCh)
				return m, nil
			}
		}
	case []list.Item:
		m.model.SetItems(msg)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm
	return m, cmd
}

func (m listModel[ItemType]) View() string {
	return m.model.View()
}

func (m listModel[ItemType]) Selected() ItemType {
	return m.model.SelectedItem().(ItemType)
}
