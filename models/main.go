package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/kubernetes"
)

type MainModel struct {
	clientset *kubernetes.Clientset
	view      ViewKey
	views     Views
	size      tea.WindowSizeMsg
	err       error
}

func Main(clientset *kubernetes.Clientset) MainModel {
	return MainModel{
		clientset: clientset,
		view:      NamespacesView,
		views:     Views{namespaces: Namespaces(clientset, defaultSize)},
		size:      defaultSize,
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.views.namespaces.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var v tea.Model
	var cmd tea.Cmd

	switch m.view {
	case NamespacesView:
		v, cmd = m.views.namespaces.Update(msg)
		m.views.namespaces = v.(*NamespacesModel)
	case ContainersView:
		v, cmd = m.views.containers.Update(msg)
		m.views.containers = v.(*ContainersModel)
	case LogsView:
		v, cmd = m.views.logs.Update(msg)
		m.views.logs = v.(*LogsModel)
	}

	if cmd != nil {
		return m, cmd
	}

	switch msg := msg.(type) {
	case ErrMsg:
		m.err = msg.Err
		return m, nil
	case tea.WindowSizeMsg:
		m.size = msg
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return MainModel{}, tea.Quit
		case tea.KeyLeft.String():
			switch m.view {
			case ContainersView:
				m.view = NamespacesView
			case LogsView:
				m.view = ContainersView
			}
		case "enter":
			switch m.view {
			case NamespacesView:
				n := m.views.namespaces
				if n == nil {
					m.err = fmt.Errorf("failed to find namespace view")
					return m, nil
				}

				namespace := n.Selected()

				m.view = ContainersView
				m.views.containers = Containers(m.clientset, m.size, namespace)

				return m, m.views.containers.Init()
			case ContainersView:
				c := m.views.containers
				if c == nil {
					m.err = fmt.Errorf("failed to find namespace view")
					return m, nil
				}

				container := c.Selected()

				m.view = LogsView
				m.views.logs = Logs(
					m.clientset,
					m.size,
					container.Namespace,
					container.Pod,
					container.Container,
				)

				return m, m.views.logs.Init()
			}
		}
	}

	return m, nil
}

func (m MainModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	switch m.view {
	case NamespacesView:
		return m.views.namespaces.View()
	case ContainersView:
		return m.views.containers.View()
	case LogsView:
		return m.views.logs.View()
	default:
		return fmt.Sprintf("unknown view: %v", m.view)
	}
}
