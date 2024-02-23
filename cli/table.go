package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableModel struct {
	table.Model
}

type tableColumn struct {
	key   string
	width int
}

type tableRow struct {
	key   string
	value string
}

func newTableModel() tableModel {
	t := table.New(
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithKeyMap(table.DefaultKeyMap()),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return tableModel{t}
}

func (m tableModel) Init() tea.Cmd { return nil }

type TableRowItem struct {
	Title string
	Value string
}

type tableRowMsg []TableRowItem

func (m tableRowMsg) columns() []table.Column {
	c := []table.Column{}

	for _, item := range m {
		c = append(c, table.Column{
			Title: item.Title,
			Width: len(fmt.Sprintf("%v", item.Title)),
		})
	}

	return c
}

func (m tableRowMsg) row() table.Row {
	r := table.Row{}

	for _, v := range m {
		r = append(r, fmt.Sprintf("%v", v))
	}

	return r
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)

	switch msg := msg.(type) {
	case tableRowMsg:
		rows := m.Model.Rows()

		if len(rows) == 0 {
			m.Model.SetColumns(msg.columns())
		}

		m.Model.SetRows(append(rows, msg.row()))
	}

	return m, cmd
}

func (m tableModel) View() string {
	return m.Model.View()
}
