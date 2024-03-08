package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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
	d := &models.ListItemDelegate{}
	d.SetShowDescription(options.showDescription)

	m := list.New([]list.Item{}, d, 0, 0)

	m.SetFilteringEnabled(true)
	m.SetShowStatusBar(false)
	m.SetSize(size.Width, size.Height)
	m.SetSpinner(spinner.Dot)
	m.Title = options.title

	// prevents esc as a quit key
	m.KeyMap.Quit = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	)

	m.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "previous page"),
			),
		}
	}

	m.Styles.NoItems = models.ListStyles.NoItems
	m.Styles.HelpStyle = models.ListStyles.Help
	m.Styles.PaginationStyle = models.ListStyles.Pagination
	m.Styles.Title = models.ListStyles.Title
	m.Styles.TitleBar = models.ListStyles.TitleBar

	return listModel[ItemType]{
		model:   &m,
		options: options,

		msgCh: msgCh,
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
