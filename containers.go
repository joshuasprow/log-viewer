package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/commands"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
	"github.com/joshuasprow/log-viewer/models"
	"k8s.io/client-go/kubernetes"
)


type containerListItem k8s.Container

func (n containerListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s/%s", n.Namespace, n.Pod, n.Name)
}

type containersModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	namespace string
	msgCh     chan<- tea.Msg
}

func newContainersModel(
	clientset *kubernetes.Clientset,
	namespace string,
	msgCh chan<- tea.Msg,
) containersModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "containers"

	return containersModel{
		clientset: clientset,
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
}

func (m containersModel) Init() tea.Cmd {
	return tea.Batch(
		m.model.StartSpinner(),
		commands.GetCronJobs(m.clientset, m.namespace),
	)
}

func (m containersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- messages.Container(m.Selected())
		}
	case messages.Containers:
		items := []list.Item{}

		for _, c := range msg {
			items = append(items, containerListItem(c))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm
	return m, cmd
}

func (m containersModel) View() string {
	return m.model.View()
}

func (m containersModel) Selected() containerListItem {
	return m.model.SelectedItem().(containerListItem)
}
