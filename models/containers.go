package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	model     list.Model
	namespace string
}

func Containers(namespace string) *ContainersModel {
	m := defaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "containers"

	m.StartSpinner()

	return &ContainersModel{
		model:     m,
		namespace: namespace,
	}
}

func (ContainersModel) Init() tea.Cmd { return nil }

func (m *ContainersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ContainersMsg:
		items := []list.Item{}

		for _, c := range msg {
			items = append(items, c)
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
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
