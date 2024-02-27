package models

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func defaultListModel() list.Model {
	m := list.New([]list.Item{}, listItemDelegate{}, 0, 0)

	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(false)

	// prevents esc as a quit key
	m.KeyMap.Quit = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	)

	m.Styles.HelpStyle = listStyles.Help
	m.Styles.PaginationStyle = listStyles.Pagination
	m.Styles.Title = listStyles.Title
	m.Styles.TitleBar = listStyles.TitleBar

	return m
}
