package models

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/styles"
)

type listItemDelegate struct{}

func (d listItemDelegate) Height() int  { return 1 }
func (d listItemDelegate) Spacing() int { return 0 }

func (d listItemDelegate) Update(
	_ tea.Msg,
	_ *list.Model,
) tea.Cmd {
	return nil
}

func (d listItemDelegate) Render(
	w io.Writer,
	m list.Model,
	index int,
	item list.Item,
) {
	i, ok := item.(list.Item)
	if !ok {
		return
	}

	fn := styles.ListItem.Normal.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.ListItem.Selected.Render("> " + strings.Join(s, " "))
		}
	}

	_, err := fmt.Fprint(w, fn(i.FilterValue()))
	if err != nil {
		// todo: panic with custom error and handle in main model
		panic(fmt.Errorf("render list item: %w", err))
	}
}
