package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
)

type jobListItem k8s.Job

func (n jobListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s", n.Namespace, n.Name)
}

type cronJobJobsModel struct {
	model   *list.Model
	cronJob k8s.CronJob
	msgCh   chan<- tea.Msg
}

func newCronJobJobsModel(
	size tea.WindowSizeMsg,
	cronJob k8s.CronJob,
	msgCh chan<- tea.Msg,
) cronJobJobsModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = "jobs"

	items := []list.Item{}

	for _, job := range cronJob.Jobs {
		items = append(items, jobListItem(job))
	}

	m.SetItems(items)

	return cronJobJobsModel{
		model:   &m,
		cronJob: cronJob,
		msgCh:   msgCh,
	}
}

func (m cronJobJobsModel) Init() tea.Cmd { return nil }

func (m cronJobJobsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.msgCh <- cronJobsViewMsg{namespace: m.cronJob.Namespace}
			return m, nil
		case "enter":
			m.msgCh <- cronJobContainersViewMsg{job: m.Selected()}
			return m, nil
		}
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm
	return m, cmd
}

func (m cronJobJobsModel) View() string {
	return m.model.View()
}

func (m cronJobJobsModel) Selected() k8s.Job {
	return k8s.Job(m.model.SelectedItem().(jobListItem))
}
