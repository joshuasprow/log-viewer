package cli

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tableDefaultStyles = table.DefaultStyles()

var tableStyles = struct {
	header   lipgloss.Style
	selected lipgloss.Style
}{
	header: tableDefaultStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false),
	selected: tableDefaultStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false),
}

type TableModel struct {
	model table.Model
}

func newTableModel() TableModel {
	t := table.New(
		table.WithFocused(true),
		table.WithKeyMap(table.DefaultKeyMap()),
	)

	t.SetStyles(tableDefaultStyles)

	return TableModel{
		model: t,
	}
}

func (m TableModel) Init() tea.Cmd { return nil }

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m TableModel) View() string {
	return m.model.View()
}

func (m TableModel) SetColumns(cols []table.Column) {
	m.model.SetColumns(cols)
}

func (m TableModel) SetRows(rows []table.Row) {
	m.model.SetRows(rows)
}

func (m TableModel) Rows() []table.Row {
	return m.model.Rows()
}

func Table() TableModel {
	return newTableModel()
}
