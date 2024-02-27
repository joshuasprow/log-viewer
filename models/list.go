package models

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func defaultListModel(size tea.WindowSizeMsg) list.Model {
	m := list.New(
		[]list.Item{},
		listItemDelegate{},
		size.Width,
		size.Height-2,
	)

	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(false)

	// prevents esc as a quit key
	m.KeyMap.Quit = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	)

	m.Styles.PaginationStyle = listStyles.Pagination
	m.Styles.HelpStyle = listStyles.Help
	m.Styles.Title = listStyles.Title
	m.Styles.TitleBar = listStyles.TitleBar

	return m
}
