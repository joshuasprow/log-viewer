package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
)

type cronJobListItem k8s.CronJob

func (n cronJobListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s", n.Namespace, n.Name)
}

type cronJobsModel struct {
	model     *list.Model
	namespace string
	msgCh     chan<- tea.Msg
}

func newCronJobsModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) cronJobsModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = "cronJobs"

	return cronJobsModel{
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
}

func (m cronJobsModel) Init() tea.Cmd {
	return m.model.StartSpinner()
}

func (m cronJobsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.msgCh <- apisViewMsg{namespace: m.namespace}
			return m, nil
		case "enter":
			m.msgCh <- cronJobJobsViewMsg{cronJob: m.Selected()}
			return m, nil
		}
	case cronJobsDataMsg:
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

func (m cronJobsModel) Selected() k8s.CronJob {
	return k8s.CronJob(m.model.SelectedItem().(cronJobListItem))
}
