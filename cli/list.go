package cli

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultListWidth = 20
const defaultListHeight = 12

var (
	listItemStyle     = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type listItem string

func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return string(i) }
func (i listItem) FilterValue() string { return string(i) }

type listModel struct {
	l list.Model
}

func newListModel(initialItems []string) listModel {
	items := []list.Item{}

	for _, i := range initialItems {
		items = append(items, listItem(i))
	}

	l := list.New(items, list.NewDefaultDelegate(), defaultListWidth, defaultListHeight)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)

	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return listModel{l}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.l, cmd = m.l.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	return m.l.View()
}
