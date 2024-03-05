package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
)

type containerListItem k8s.Container

func (n containerListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s/%s", n.Namespace, n.Pod, n.Name)
}

type containersModel struct {
	model     *list.Model
	namespace string
	msgCh     chan<- tea.Msg
}

func newContainersModel(
	size tea.WindowSizeMsg,
	namespace string,
	msgCh chan<- tea.Msg,
) containersModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.SetSize(size.Width, size.Height)
	m.Title = "containers"

	return containersModel{
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
}

func (m containersModel) Init() tea.Cmd {
	return tea.Batch(m.model.StartSpinner())
}

func (m containersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.msgCh <- viewMsg{
				key:  apisKey,
				data: m.namespace,
			}
		case "enter":
			m.msgCh <- viewMsg{
				key:  containerLogsKey,
				data: k8s.Container(m.Selected()),
			}
		}
	case containersDataMsg:
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

func (m containersModel) View() string {
	return m.model.View()
}

func (m containersModel) Selected() containerListItem {
	return m.model.SelectedItem().(containerListItem)
}
