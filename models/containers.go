package models

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"k8s.io/client-go/kubernetes"
)

type ContainerListItem struct {
	Namespace string
	Pod       string
	Container string
}

func (n ContainerListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s/%s", n.Namespace, n.Pod, n.Container)
}

type ContainersModel struct {
	clientset *kubernetes.Clientset
	model     list.Model
	namespace string
}

func Containers(
	clientset *kubernetes.Clientset,
	size tea.WindowSizeMsg,
	namespace string,
) *ContainersModel {
	m := defaultListModel(size)
	m.SetFilteringEnabled(true)
	m.Title = "containers"

	return &ContainersModel{
		clientset: clientset,
		model:     m,
		namespace: namespace,
	}
}

func (m *ContainersModel) initData() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		pods, err := k8s.GetPods(ctx, m.clientset, m.namespace)
		if err != nil {
			return ErrMsg{Err: fmt.Errorf("load model data: %w", err)}
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

		return ContainersMsg(containers)
	}
}

func (m *ContainersModel) Init() tea.Cmd {
	return tea.Batch(m.model.StartSpinner(), m.initData())
}

func (m *ContainersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrMsg:
		m.model.StopSpinner()
		return m, nil // todo: return error Cmd ?
	case ContainersMsg:
		items := []list.Item{}

		for _, c := range msg {
			items = append(items, c)
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	case tea.WindowSizeMsg:
		m.model.SetWidth(msg.Width)
		m.model.SetHeight(msg.Height - 2)
	}

	var cmd tea.Cmd

	m.model, cmd = m.model.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

func (m *ContainersModel) View() string {
	return m.model.View()
}

func (m *ContainersModel) Selected() ContainerListItem {
	return m.model.SelectedItem().(ContainerListItem)
}
