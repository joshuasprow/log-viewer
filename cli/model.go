package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshuasprow/log-viewer/pkg"
)

type viewKey string

const (
	namespaceKey viewKey = "namespace"
	podKey       viewKey = "pod"
	logsKey      viewKey = "logs"
)

type model struct {
	quitting       bool
	view           viewKey
	namespaceModel tea.Model
	namespace      string
	podModel       tea.Model
	pod            string
	logsModel      tea.Model
	logCh          <-chan TableRowItem
	logModel       tea.Model
	log            string
	err            error
}

func NewModel(rowCh <-chan TableRowItem) tea.Model {
	return model{
		view:           namespaceKey,
		namespaceModel: newListModel(pkg.Namespaces),
		podModel:       newListModel([]string{}),
		logsModel:      newTableModel(rowCh),
		logCh:          rowCh,
	}
}

func waitForLog(rowCh <-chan TableRowItem) tea.Cmd {
	return func() tea.Msg {
		row := <-rowCh

		if row == (TableRowItem{}) {
			return waitForLog(rowCh)
		}

		return row
	}
}

func (m model) Init() tea.Cmd {
	return waitForLog(m.logCh)
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

	m.namespaceModel, cmd = m.namespaceModel.Update(msg)
	m.podModel, cmd = m.podModel.Update(msg)
	m.logsModel, cmd = m.logsModel.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true

			return m, tea.Quit
		case "enter":
			switch m.view {
			case namespaceKey:
				s, err := getSelectedListItem(m.namespaceModel)
				if err != nil {
					m.err = fmt.Errorf("failed to get selected namespace: %w", err)
					return m, nil
				}

				m.namespace = s

				pods, ok := pkg.Pods[m.namespace]
				if !ok {
					m.err = fmt.Errorf("no pods found in namespace %s", m.namespace)
					return m, nil
				}

				m.podModel = newListModel(pods)

				m.view = podKey
			case podKey:
				r, err := getSelectedListItem(m.podModel)
				if err != nil {
					m.err = fmt.Errorf("failed to get selected pod: %w", err)
					return m, nil
				}

				m.pod = r
				m.logsModel = newTableModel(m.logCh)

				m.view = logsKey
			}
		case "esc":
			switch m.view {
			case podKey:
				m.pod = ""
				m.podModel = newListModel([]string{})

				m.view = namespaceKey
			case logsKey:
				m.log = ""
				m.logsModel = newTableModel(m.logCh)

				m.view = podKey
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
	case namespaceKey:
		body = m.namespaceModel.View()
	case podKey:
		body = m.podModel.View()
	case logsKey:
		body = m.logsModel.View()
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
			titlePart("source", m.namespace),
			titlePart("resource", m.pod),
		),
		body,
	)
}
