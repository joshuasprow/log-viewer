package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
)

type appModel struct {
	msgCh chan<- tea.Msg
	size  tea.WindowSizeMsg
	view  tea.Model
}

func newAppModel(msgCh chan<- tea.Msg) appModel {
	defaultSize := tea.WindowSizeMsg{Width: 80, Height: 24}

	return appModel{
		msgCh: msgCh,
		size:  defaultSize,
		view:  newNamespacesModel(defaultSize, msgCh),
	}
}

func (m appModel) Init() tea.Cmd {
	return m.view.Init()
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.Width = msg.Width
		m.size.Height = msg.Height - 1 // todo: fixes list title disappearing
	case messages.Namespace:
		switch msg.Api {
		case "":
			m.view = newApisViewModel(m.size, msg.Name, m.msgCh)
			return m, m.view.Init()
		case messages.ContainersApi:
			m.view = newContainersModel(m.size, msg.Name, m.msgCh)
			return m, m.view.Init()
		case messages.CronJobsApi:
			m.view = newCronJobsModel(m.size, msg.Name, m.msgCh)
			return m, m.view.Init()
		default:
			panic(fmt.Errorf("unknown namespace view: %s", msg.Api))
		}
	case messages.Container:
		m.view = newLogsModel(m.size, k8s.Container(msg), m.msgCh)
		return m, m.view.Init()
	case messages.CronJob:
		m.view = newJobsModel(m.size, msg.Jobs, m.msgCh)
		return m, m.view.Init()
	case messages.Job:
		m.view = newJobContainersModel(m.size, k8s.Job(msg), m.msgCh)
		return m, m.view.Init()
	case messages.JobContainer:
		m.view = newLogsModel(m.size, k8s.Container(msg), m.msgCh)
		return m, m.view.Init()
	}

	var cmd tea.Cmd

	if m.view != nil {
		m.view, cmd = m.view.Update(msg)
	}

	return m, cmd
}

func (m appModel) View() string {
	var v string
	if m.view == nil {
		v = spinner.New().View()
	} else {
		v = m.view.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, v)
}
