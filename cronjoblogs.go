package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
)

type cronJobLogsModel struct {
	model *list.Model
	msgCh chan<- tea.Msg
	job   k8s.Job
}

func newCronJobLogsModel(
	size tea.WindowSizeMsg,
	job k8s.Job,
	msgCh chan<- tea.Msg,
) cronJobLogsModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = "cronJob logs"

	return cronJobLogsModel{
		model: &m,
		msgCh: msgCh,
		job:   job,
	}
}

func (m cronJobLogsModel) Init() tea.Cmd {
	return m.model.StartSpinner()
}

func (m cronJobLogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.msgCh <- viewMsg{
				key:  cronJobContainersKey,
				data: m.job,
			}
		}
	case cronJobLogsDataMsg:
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

func (m cronJobLogsModel) View() string {
	return m.model.View()
}

func (m cronJobLogsModel) Selected() list.Item {
	return m.model.SelectedItem()
}
