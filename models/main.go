package models

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
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
		views:     Views{namespaces: Namespaces(defaultSize)},
		size:      defaultSize,
	}
}

func (m MainModel) Init() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		namespaces, err := k8s.GetNamespaces(ctx, m.clientset)
		if err != nil {
			return ErrMsg{Err: fmt.Errorf("load model data: %w", err)}
		}

		return NamespacesMsg(namespaces)
	}
}

func getContainers(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
) (
	[]ContainerListItem,
	error,
) {
	pods, err := k8s.GetPods(ctx, clientset, namespace)
	if err != nil {
		return nil, fmt.Errorf("load model data: %w", err)
	}

	containers := []ContainerListItem{}

	for _, pod := range pods {
		if len(pod.Spec.Containers) == 0 {
			containers = append(containers, ContainerListItem{
				Namespace: pod.Namespace,
				Pod:       pod.Name,
			})
			continue
		}

		for _, container := range pod.Spec.Containers {
			containers = append(containers, ContainerListItem{
				Namespace: pod.Namespace,
				Pod:       pod.Name,
				Container: container.Name,
			})
		}
	}

	return containers, nil
}

func (m MainModel) handleEnter() (MainModel, tea.Cmd) {
	switch m.view {
	case NamespacesView:
		n := m.views.namespaces
		if n == nil {
			m.err = fmt.Errorf("failed to find namespace view")
			return m, nil
		}

		namespace := n.Selected()

		m.view = ContainersView
		m.views.containers = Containers(m.size, namespace)

		return m, func() tea.Msg {
			ctx := context.Background()

			containers, err := getContainers(ctx, m.clientset, namespace)
			if err != nil {
				return ErrMsg{Err: fmt.Errorf("get containers: %w", err)}
			}

			return ContainersMsg(containers)
		}
	case ContainersView:
		c := m.views.containers
		if c == nil {
			m.err = fmt.Errorf("failed to find containers view")
			return m, nil
		}

		container := c.Selected()

		m.view = LogsView
		m.views.logs = Logs(m.size, container)

		return m, func() tea.Msg {
			ctx := context.Background()

			logs, err := k8s.GetPodLogs(
				ctx,
				m.clientset,
				container.Namespace,
				container.Pod,
				container.Container,
			)
			if err != nil {
				return ErrMsg{Err: fmt.Errorf("get pod logs: %w", err)}
			}

			return LogsMsg(logs)
		}
	}

	return m, nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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
			m, cmd = m.handleEnter()
			if cmd != nil {
				return m, cmd
			}
		}
	}

	var v tea.Model

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

	return m, cmd
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
