package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/messages"
	"k8s.io/client-go/kubernetes"
)

type namespaceListItem string

func (n namespaceListItem) FilterValue() string {
	return string(n)
}

type NamespacesModel struct {
	clientset *kubernetes.Clientset
	// a pointer is necessary for updating the spinner state
	model *list.Model
}

func Namespaces() NamespacesModel {
	m := DefaultListModel()
	m.Title = "namespaces"

	return NamespacesModel{model: &m}
}

func (NamespacesModel) Init() tea.Cmd { return nil }

func (m NamespacesModel) Update(msg tea.Msg) (NamespacesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.Namespaces:
		items := []list.Item{}

		for _, namespace := range msg {
			items = append(items, namespaceListItem(namespace))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm

	return m, cmd
}

func (m NamespacesModel) View() string {
	return m.model.View()
}

func (m NamespacesModel) Selected() string {
	return m.model.SelectedItem().(namespaceListItem).FilterValue()
}
