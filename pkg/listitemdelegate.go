package pkg

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type listItemDelegate struct{}

func NewListItemDelegate() list.ItemDelegate {
	return listItemDelegate{}
}

func (d listItemDelegate) Height() int                             { return 1 }
func (d listItemDelegate) Spacing() int                            { return 0 }
func (d listItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(list.Item)
	if !ok {
		return
	}

	fn := ListItemStyles.Normal.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return ListItemStyles.Selected.Render("> " + strings.Join(s, " "))
		}
	}

	_, err := fmt.Fprint(w, fn(i.FilterValue()))
	if err != nil {
		// todo: panic with custom error and handle in main model
		panic(fmt.Errorf("render list item: %w", err))
	}
}
