package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/kubernetes"
)

type namespaceListItem string

func (n namespaceListItem) FilterValue() string {
	return string(n)
}

type NamespacesModel struct {
	clientset *kubernetes.Clientset
	model     list.Model
}

func Namespaces() *NamespacesModel {
	m := defaultListModel()
	m.Title = "namespaces"

	m.StartSpinner()

	return &NamespacesModel{model: m}
}

func (NamespacesModel) Init() tea.Cmd { return nil }

func (m *NamespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NamespacesMsg:
		items := []list.Item{}

		for _, namespace := range msg {
			items = append(items, namespaceListItem(namespace))
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

func (m *NamespacesModel) View() string {
	return m.model.View()
}

func (m *NamespacesModel) Selected() string {
	return m.model.SelectedItem().(namespaceListItem).FilterValue()
}
