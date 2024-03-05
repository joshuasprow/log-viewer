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

type jobsModel struct {
	model     *list.Model
	namespace string
	msgCh     chan<- tea.Msg
}

func newJobsModel(
	size tea.WindowSizeMsg,
	namespace string,
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
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
}

func (m jobsModel) Init() tea.Cmd { return nil }

func (m jobsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.msgCh <- viewMsg{
				key:  cronJobsKey,
				data: m.namespace,
			}
		case "enter":
			m.msgCh <- viewMsg{
				key:  cronJobContainersKey,
				data: k8s.Job(m.Selected()),
			}
		}
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
