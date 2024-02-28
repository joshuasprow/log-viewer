package models

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
	"k8s.io/client-go/kubernetes"
)

type MainModel struct {
	clientset *kubernetes.Clientset
	view      ViewKey
	size      tea.WindowSizeMsg
	err       error

	// views
	namespaces NamespacesModel
	containers ContainersModel
	logs       LogsModel
}

func Main(clientset *kubernetes.Clientset) MainModel {
	return MainModel{
		clientset: clientset,
		view:      NamespacesView,
		size: tea.WindowSizeMsg{
			Width:  appStyles.GetWidth(),
			Height: appStyles.GetHeight(),
		},

		namespaces: Namespaces(),
		containers: Containers(),
		logs:       Logs(),
	}
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		m.namespaces.model.StartSpinner(),
		func() tea.Msg {
			ctx := context.Background()

			namespaces, err := k8s.GetNamespaces(ctx, m.clientset)
			if err != nil {
				return messages.Error{Err: fmt.Errorf("load model data: %w", err)}
			}

			return messages.Namespaces(namespaces)
		},
	)
}

func (m MainModel) handleEnter() (MainModel, tea.Cmd) {
	switch m.view {
	case NamespacesView:
		namespace := m.namespaces.Selected()

		m.view = ContainersView
		m.containers = Containers()

		return m, tea.Batch(
			m.containers.model.StartSpinner(),
			func() tea.Msg {
				ctx := context.Background()

				containers, err := k8s.GetContainers(ctx, m.clientset, namespace)
				if err != nil {
					return messages.Error{Err: fmt.Errorf("get containers: %w", err)}
				}

				return messages.Containers(containers)
			},
		)
	case ContainersView:
		container := m.containers.Selected()

		m.view = LogsView
		m.logs = Logs()

		return m, tea.Batch(
			m.logs.model.StartSpinner(),
			func() tea.Msg {
				ctx := context.Background()

				logs, err := k8s.GetPodLogs(
					ctx,
					m.clientset,
					container.Namespace,
					container.Pod,
					container.Name,
				)
				if err != nil {
					return messages.Error{Err: fmt.Errorf("get pod logs: %w", err)}
				}

				return messages.Logs(logs)
			},
		)
	}

	return m, nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case messages.Error:
		m.err = msg.Err
		return m, nil
	case tea.WindowSizeMsg:
		// todo: this can't be right
		x, y := appStyles.GetFrameSize()

		m.size = tea.WindowSizeMsg{
			Width:  msg.Width - x,
			Height: msg.Height - y,
		}
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return MainModel{}, tea.Quit
		case "esc":
			switch m.view {
			case ContainersView:
				m.view = NamespacesView
			case LogsView:
				m.view = ContainersView
			}
		case "enter":
			m, cmd = m.handleEnter()
			if cmd != nil {
				return m, cmd
			}
		}
	}

	switch m.view {
	case NamespacesView:
		m.namespaces, cmd = m.namespaces.Update(msg)
		m.namespaces.model.SetSize(m.size.Width, m.size.Height)
	case ContainersView:
		m.containers, cmd = m.containers.Update(msg)
		m.containers.model.SetSize(m.size.Width, m.size.Height)
	case LogsView:
		m.logs, cmd = m.logs.Update(msg)
		m.logs.model.SetSize(m.size.Width, m.size.Height)
	}

	return m, cmd
}

func (m MainModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	switch m.view {
	case NamespacesView:
		return m.namespaces.View()
	case ContainersView:
		return m.containers.View()
	case LogsView:
		return m.logs.View()
	default:
		return fmt.Sprintf("unknown view: %v", m.view)
	}
}
