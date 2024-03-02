package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/commands"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
	"github.com/joshuasprow/log-viewer/models"
	"k8s.io/client-go/kubernetes"
)

type logListItem string

func (i logListItem) FilterValue() string {
	return string(i)
}

type logsModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	msgCh     chan<- tea.Msg

	namespace string
	pod       string
	container string
}

func newLogsModel(
	clientset *kubernetes.Clientset,
	container k8s.Container,
	msgCh chan<- tea.Msg,
) logsModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "logs"

	return logsModel{
		clientset: clientset,
		model:     &m,
		msgCh:     msgCh,

		namespace: container.Namespace,
		pod:       container.Pod,
		container: container.Name,
	}
}

func (m logsModel) Init() tea.Cmd {
	return tea.Batch(
		m.model.StartSpinner(),
		commands.GetLogs(m.clientset, m.namespace, m.pod, m.container),
	)
}

func (m logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.Logs:
		items := []list.Item{}

		for _, i := range msg {
			items = append(items, logListItem(i))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm

	return m, cmd
}

func (m logsModel) View() string {
	return m.model.View()
}

func (m logsModel) Selected() list.Item {
	return m.model.SelectedItem()
}
