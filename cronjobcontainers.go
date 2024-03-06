package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
)

type cronJobContainersModel struct {
	model   *list.Model
	cronJob k8s.CronJob
	msgCh   chan<- tea.Msg
}

func newCronJobContainersModel(
	size tea.WindowSizeMsg,
	cronJob k8s.CronJob,
	msgCh chan<- tea.Msg,
) cronJobContainersModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = renderTitle(cronJob.Namespace, cronJob.Name, "select a container")

	return cronJobContainersModel{
		model:   &m,
		cronJob: cronJob,
		msgCh:   msgCh,
	}
}

func (m cronJobContainersModel) Init() tea.Cmd {
	return m.model.StartSpinner()
}

func (m cronJobContainersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.msgCh <- cronJobJobsViewMsg{cronJob: m.cronJob}
			return m, nil
		case "enter":
			m.msgCh <- cronJobLogsViewMsg{container: m.Selected()}
			return m, nil
		}
	case cronJobContainersDataMsg:
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

func (m cronJobContainersModel) View() string {
	return m.model.View()
}

func (m cronJobContainersModel) Selected() k8s.Container {
	return k8s.Container(m.model.SelectedItem().(containerListItem))
}
