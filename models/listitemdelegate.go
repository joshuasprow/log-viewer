package models

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItemDelegate struct {
	showDescription bool
	height          int
	width           int
}

func (d *ListItemDelegate) SetShowDescription(b bool) {
	d.showDescription = b
	if b {
		d.height = 2
	}
}

func (d ListItemDelegate) Height() int {
	if d.showDescription {
		return d.height
	}
	return 1
}

func (d *ListItemDelegate) SetHeight(h int) {
	d.height = h
}

func (d ListItemDelegate) Spacing() int { return 0 }

func (d *ListItemDelegate) Update(
	_ tea.Msg,
	m *list.Model,
) tea.Cmd {
	d.width = m.Width()
	if d.height == 0 {
		d.height = 1
	}
	return nil
}

func (d *ListItemDelegate) Render(
	w io.Writer,
	m list.Model,
	index int,
	item list.Item,
) {
	var title string

	if ti, ok := item.(Titled); ok {
		title = ti.Title()
	} else {
		title = item.FilterValue()
	}

	if index == m.Index() {
		title = listItemStyles.SelectedTitle.Render("> " + title)
	} else {
		title = listItemStyles.NormalTitle.Render(title)
	}

	var desc string

	if di, ok := item.(Described); d.showDescription && ok {
		desc = listItemStyles.Description.Render(di.Description())
	}

	if desc == "" {
		fmt.Fprint(w, title)
		return
	}

	fmt.Fprint(w, title+"\n"+desc)
}
