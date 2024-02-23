package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const defaultListWidth = 20
const defaultListHeight = 12

var listItemStyles = struct {
	normal   lipgloss.Style
	selected lipgloss.Style
}{
	normal:   lipgloss.NewStyle().PaddingLeft(4),
	selected: lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")),
}

var listStyles = struct {
	pagination lipgloss.Style
	help       lipgloss.Style
	quitText   lipgloss.Style
}{
	pagination: lipgloss.NewStyle().PaddingLeft(4),
	help:       list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1),
	quitText:   lipgloss.NewStyle().Margin(1, 0, 2, 4),
}

type listItem string

func (i listItem) FilterValue() string { return string(i) }

type listItemDelegate struct{}

func (d listItemDelegate) Height() int                             { return 1 }
func (d listItemDelegate) Spacing() int                            { return 0 }
func (d listItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(listItem)
	if !ok {
		return
	}

	fn := listItemStyles.normal.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return listItemStyles.normal.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(i.FilterValue()))
}

type listModel struct {
	l list.Model
}

func newListModel(initialItems []string) listModel {
	items := []list.Item{}

	for _, i := range initialItems {
		items = append(items, listItem(i))
	}

	l := list.New(items, listItemDelegate{}, defaultListWidth, defaultListHeight)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)

	l.Styles.PaginationStyle = listStyles.pagination
	l.Styles.HelpStyle = listStyles.help

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
