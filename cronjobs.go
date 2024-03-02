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


type cronJobListItem k8s.Container

func (n cronJobListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s/%s", n.Namespace, n.Pod, n.Name)
}

type cronJobsModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	namespace string
	msgCh     chan<- tea.Msg
}

func newCronJobsModel(
	clientset *kubernetes.Clientset,
	namespace string,
	msgCh chan<- tea.Msg,
) cronJobsModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "cronJobs"

	return cronJobsModel{
		clientset: clientset,
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
}

func (m cronJobsModel) Init() tea.Cmd {
	return tea.Batch(
		m.model.StartSpinner(),
		commands.GetContainers(m.clientset, m.namespace),
	)
}

func (m cronJobsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- messages.Container(m.Selected())
		}
	case messages.Containers:
		items := []list.Item{}

		for _, c := range msg {
			items = append(items, cronJobListItem(c))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm
	return m, cmd
}

func (m cronJobsModel) View() string {
	return m.model.View()
}

func (m cronJobsModel) Selected() cronJobListItem {
	return m.model.SelectedItem().(cronJobListItem)
}
