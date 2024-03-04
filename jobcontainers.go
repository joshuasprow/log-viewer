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

type jobContainerListItem k8s.Container

func (n jobContainerListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s", n.Namespace, n.Name)
}

type jobContainersModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	job       k8s.Job
	msgCh     chan<- tea.Msg
}

func newJobContainersModel(
	clientset *kubernetes.Clientset,
	size tea.WindowSizeMsg,
	job k8s.Job,
	msgCh chan<- tea.Msg,
) jobContainersModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = "jobContainers"

	return jobContainersModel{
		clientset: clientset,
		model:     &m,
		job:       job,
		msgCh:     msgCh,
	}
}

func (m jobContainersModel) Init() tea.Cmd {
	return tea.Batch(
		m.model.StartSpinner(),
		commands.GetJobContainers(m.clientset, m.job.Namespace, m.job.Name),
	)
}

func (m jobContainersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- messages.JobContainer(m.Selected())
		}
	case messages.JobContainers:
		items := []list.Item{}

		for _, c := range msg {
			items = append(items, jobContainerListItem(c))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm
	return m, cmd
}

func (m jobContainersModel) View() string {
	return m.model.View()
}

func (m jobContainersModel) Selected() jobContainerListItem {
	return m.model.SelectedItem().(jobContainerListItem)
}
