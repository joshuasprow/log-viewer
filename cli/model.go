package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewKey string

const (
	sourceKey   viewKey = "source"
	resourceKey viewKey = "resource"
	viewerKey   viewKey = "viewer"
)

var views = map[viewKey]tea.Model{
	sourceKey:   newListModel([]string{"s-1", "s-2", "s-3"}),
	resourceKey: newListModel([]string{"r-1", "r-2", "r-3"}),
	viewerKey:   newTableModel(),
}

type model struct {
	quitting      bool
	view          viewKey
	sourceModel   tea.Model
	source        string
	resourceModel tea.Model
	resource      string
	viewerModel   tea.Model
	err           error
}

func NewModel(rowCh chan<- table.Row) tea.Model {
	return model{
		view:          sourceKey,
		sourceModel:   newListModel([]string{"s-1", "s-2", "s-3"}),
		resourceModel: newListModel([]string{"r-1", "r-2", "r-3"}),
		viewerModel: newTableModel(
			[]table.Column{
				{Title: "Rank", Width: 4},
				{Title: "City", Width: 10},
				{Title: "Country", Width: 10},
				{Title: "Population", Width: 10},
			},
			[]table.Row{},
		),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func getSelectedListItem(m tea.Model) (string, error) {
	lm, ok := m.(listModel)
	if !ok {
		return "", fmt.Errorf("expected model to be a listModel, but got %T", m)
	}

	i, ok := lm.l.SelectedItem().(listItem)
	if !ok {
		return "", fmt.Errorf("expected listModel.SelectedItem to return an item, but got %T", lm.l.SelectedItem())
	}

	return string(i), nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.sourceModel, cmd = m.sourceModel.Update(msg)
	m.resourceModel, cmd = m.resourceModel.Update(msg)
	m.viewerModel, cmd = m.viewerModel.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true

			return m, tea.Quit
		case "enter":
			switch m.view {
			case sourceKey:
				s, err := getSelectedListItem(m.sourceModel)
				if err != nil {
					m.err = err
					return m, nil
				}

				m.source = s
				m.view = resourceKey
			case resourceKey:
				r, err := getSelectedListItem(m.resourceModel)
				if err != nil {
					m.err = err
				} else {
					m.resource = r
					m.view = viewerKey
				}
			}
		case "esc":
			switch m.view {
			case resourceKey:
				m.view = sourceKey
				m.resource = ""
			case viewerKey:
				m.view = resourceKey
				m.viewerModel = newTableModel([]table.Column{}, []table.Row{})
			}
		}
	}

	if m.err != nil {
		panic(m.err)
	}

	return m, cmd
}

func (m model) View() string {
	body := ""

	switch m.view {
	case sourceKey:
		body = m.sourceModel.View()
	case resourceKey:
		body = m.resourceModel.View()
	case viewerKey:
		body = m.viewerModel.View()
	default:
		m.err = fmt.Errorf("unknown view key: %q", m.view)

		return m.err.Error()
	}

	titlePart := func(k, v string) string {
		return lipgloss.NewStyle().PaddingLeft(2).Render(fmt.Sprintf("%s: %q", k, v))
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			titlePart("view", string(m.view)),
			titlePart("source", m.source),
			titlePart("resource", m.resource),
		),
		body,
	)
}
