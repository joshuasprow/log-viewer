package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/messages"
	"k8s.io/client-go/kubernetes"
)

type cronJobListItem string

func (n cronJobListItem) FilterValue() string {
	return string(n)
}

type CronJobsModel struct {
	clientset *kubernetes.Clientset
	// a pointer is necessary for updating the spinner state
	model *list.Model
}

func CronJobs() CronJobsModel {
	m := DefaultListModel()
	m.Title = "cronJobs"

	return CronJobsModel{model: &m}
}

func (CronJobsModel) Init() tea.Cmd { return nil }

func (m CronJobsModel) Update(msg tea.Msg) (CronJobsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.CronJobs:
		items := []list.Item{}

		for _, cronJob := range msg {
			items = append(items, cronJobListItem(cronJob.Name))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm

	return m, cmd
}

func (m CronJobsModel) View() string {
	return m.model.View()
}

func (m CronJobsModel) Selected() string {
	return m.model.SelectedItem().(cronJobListItem).FilterValue()
}
