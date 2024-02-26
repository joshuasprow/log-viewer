package models

import (
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

	m.Styles.PaginationStyle = listStyles.Pagination
	m.Styles.HelpStyle = listStyles.Help
	m.Styles.Title = listStyles.Title
	m.Styles.TitleBar = listStyles.TitleBar

	return m
}
