package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
)

type ContainerListItem k8s.Container

func (n ContainerListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s/%s", n.Namespace, n.Pod, n.Name)
}

type ContainersModel struct {
	model list.Model
}

func Containers() *ContainersModel {
	m := defaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "containers"

	return &ContainersModel{model: m}
}

func (ContainersModel) Init() tea.Cmd { return nil }

func (m *ContainersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ContainersMsg:
		items := []list.Item{}

		for _, c := range msg {
			items = append(items, ContainerListItem(c))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	var cmd tea.Cmd

	m.model, cmd = m.model.Update(msg)

	return m, cmd
}

func (m *ContainersModel) View() string {
	return m.model.View()
}

func (m *ContainersModel) Selected() ContainerListItem {
	return m.model.SelectedItem().(ContainerListItem)
}
