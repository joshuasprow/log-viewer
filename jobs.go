package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
	"github.com/joshuasprow/log-viewer/models"
	"k8s.io/client-go/kubernetes"
)

type jobListItem k8s.Job

func (n jobListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s", n.Namespace, n.Name)
}

type jobsModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	msgCh     chan<- tea.Msg
}

func newJobsModel(
	clientset *kubernetes.Clientset,
	size tea.WindowSizeMsg,
	jobs []k8s.Job,
	msgCh chan<- tea.Msg,
) jobsModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = "jobs"

	items := []list.Item{}

	for _, job := range jobs {
		items = append(items, jobListItem(job))
	}

	m.SetItems(items)

	return jobsModel{
		clientset: clientset,
		model:     &m,
		msgCh:     msgCh,
	}
}

func (m jobsModel) Init() tea.Cmd {
	return nil
}

func (m jobsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- messages.Job(m.Selected())
		}
	case messages.Jobs:
		items := []list.Item{}

		for _, c := range msg {
			items = append(items, jobListItem(c))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm
	return m, cmd
}

func (m jobsModel) View() string {
	return m.model.View()
}

func (m jobsModel) Selected() jobListItem {
	return m.model.SelectedItem().(jobListItem)
}
