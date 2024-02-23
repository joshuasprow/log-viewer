package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableModel struct {
	model table.Model
	cols  []table.Column
	rowCh <-chan TableRowItem
}

type TableRowItem struct {
	Level string
	Time  string
	Msg   string
	Raw   string
}

func newTableModel(rowCh <-chan TableRowItem) tableModel {
	cols := []table.Column{}

	for _, title := range []string{"level", "time", "message"} {
		cols = append(cols, table.Column{
			Title: title,
			Width: len(title),
		})
	}

	t := table.New(
		table.WithColumns(cols),
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

	return tableModel{
		model: t,
		cols:  cols,
		rowCh: rowCh,
	}
}

func waitForRowMsg(rowCh <-chan TableRowItem) tea.Cmd {
	return func() tea.Msg {
		row := <-rowCh

		if row == (TableRowItem{}) {
			return nil
		}

		return row
	}
}

func (m tableModel) Init() tea.Cmd {
	return waitForRowMsg(m.rowCh)
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case TableRowItem:
		rows := m.model.Rows()
		row := []string{msg.Level, msg.Time, msg.Msg}

		for i := range m.cols {
			if len(row[i]) > m.cols[i].Width {
				m.cols[i].Width = len(row[i])
			}
		}

		m.model.SetColumns(m.cols)
		m.model.SetRows(append(rows, row))

		m.model.GotoBottom()

		return m, waitForRowMsg(m.rowCh)
	}

	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)

	return m, cmd
}

func (m tableModel) View() string {
	return m.model.View()
}
