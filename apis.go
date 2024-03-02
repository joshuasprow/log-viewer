package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/messages"
	"github.com/joshuasprow/log-viewer/models"
	"k8s.io/client-go/kubernetes"
)

type apisListItem string

func (n apisListItem) FilterValue() string {
	return string(n)
}

type apisModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	namespace string
	msgCh     chan<- tea.Msg
}

func newApisViewModel(
	clientset *kubernetes.Clientset,
	namespace string,
	msgCh chan<- tea.Msg,
) apisModel {
	m := models.DefaultListModel()
	m.Title = "apis"
	m.SetItems([]list.Item{
		apisListItem(messages.ContainersApi),
		apisListItem(messages.CronJobsApi),
	})

	return apisModel{
		clientset: clientset,
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
}

func (m apisModel) Init() tea.Cmd {
	return nil
}

func (m apisModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- messages.Namespace{
				Name: m.namespace,
				Api:  messages.Api(m.Selected()),
			}
		}
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm

	return m, cmd
}

func (m apisModel) View() string {
	return m.model.View()
}

func (m apisModel) Selected() apisListItem {
	return m.model.SelectedItem().(apisListItem)
}
