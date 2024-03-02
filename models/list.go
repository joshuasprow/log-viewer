package models

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/joshuasprow/log-viewer/styles"
)

func DefaultListModel() list.Model {
	m := list.New([]list.Item{}, listItemDelegate{}, 0, 0)

	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(false)
	m.SetSpinner(spinner.Dot)

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

	m.Styles.NoItems = styles.List.NoItems
	m.Styles.HelpStyle = styles.List.Help
	m.Styles.PaginationStyle = styles.List.Pagination
	m.Styles.Title = styles.List.Title
	m.Styles.TitleBar = styles.List.TitleBar

	return m
}
