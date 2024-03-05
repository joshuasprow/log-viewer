package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
)

type jobContainerListItem k8s.Container

func (n jobContainerListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s", n.Namespace, n.Name)
}

type jobContainersModel struct {
	model   *list.Model
	cronJob k8s.CronJob
	msgCh   chan<- tea.Msg
}

func newJobContainersModel(
	size tea.WindowSizeMsg,
	cronJob k8s.CronJob,
	msgCh chan<- tea.Msg,
) jobContainersModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = "jobContainers"

	return jobContainersModel{
		model:   &m,
		cronJob: cronJob,
		msgCh:   msgCh,
	}
}

func (m jobContainersModel) Init() tea.Cmd {
	return m.model.StartSpinner()
}

func (m jobContainersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.msgCh <- viewMsg{
				key:  cronJobJobsKey,
				data: m.cronJob,
			}
		case "enter":
			m.msgCh <- viewMsg{
				key:  cronJobLogsKey,
				data: k8s.Container(m.Selected()),
			}
		}
	case cronJobContainersDataMsg:
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
