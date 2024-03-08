package defaults

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type ListModel[ItemType any] struct {
	model   *list.Model
	options ListModelOptions[ItemType]
	msgCh   chan<- tea.Msg
}

type ListModelOptions[ItemType any] struct {
	OnEnter         func(selected ItemType, msgCh chan<- tea.Msg)
	OnEsc           func(msgCh chan<- tea.Msg)
	ShowDescription bool
	Title           string
}

func NewListModel[ItemType any](
	size tea.WindowSizeMsg,
	options ListModelOptions[ItemType],
	msgCh chan<- tea.Msg,
) ListModel[ItemType] {
	d := &ListItemDelegate{}
	d.SetShowDescription(options.ShowDescription)

	m := list.New([]list.Item{}, d, 0, 0)

	m.SetFilteringEnabled(true)
	m.SetShowStatusBar(false)
	m.SetSize(size.Width, size.Height)
	m.SetSpinner(spinner.Dot)
	m.Title = options.Title

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

	m.Styles.NoItems = ListStyles.NoItems
	m.Styles.HelpStyle = ListStyles.Help
	m.Styles.PaginationStyle = ListStyles.Pagination
	m.Styles.Title = ListStyles.Title
	m.Styles.TitleBar = ListStyles.TitleBar

	return ListModel[ItemType]{
		model:   &m,
		options: options,

		msgCh: msgCh,
	}
}

func (m ListModel[ItemType]) Init() tea.Cmd {
	return m.model.StartSpinner()
}

func (m ListModel[ItemType]) Update(
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
			if m.options.OnEsc != nil {
				m.options.OnEsc(m.msgCh)
				return m, nil
			}
		case "enter":
			if m.options.OnEnter != nil {
				m.options.OnEnter(m.Selected(), m.msgCh)
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

func (m ListModel[ItemType]) View() string {
	return m.model.View()
}

func (m ListModel[ItemType]) Selected() ItemType {
	return m.model.SelectedItem().(ItemType)
}
