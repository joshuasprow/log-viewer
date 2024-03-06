package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appModel struct {
	msgCh chan<- tea.Msg
	size  tea.WindowSizeMsg
	view  tea.Model
	data  viewData
}

func newAppModel(msgCh chan<- tea.Msg) appModel {
	size := tea.WindowSizeMsg{Width: 80, Height: 24}

	return appModel{
		msgCh: msgCh,
		size:  size,
	}
}

func (m appModel) Init() tea.Cmd {
	return func() tea.Msg {
		m.msgCh <- namespacesViewMsg{}
		return nil
	}
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.Width = msg.Width
		m.size.Height = msg.Height - 1 // todo: fixes list title disappearing
	case namespacesViewMsg:
		m.view = newNamespacesModel(m.size, m.msgCh)

		return m, m.view.Init()
	case apisViewMsg:
		m.data.namespace = msg.namespace
		m.view = newApisModel(m.size, m.data.namespace, m.msgCh)

		return m, m.view.Init()
	case containersViewMsg:
		m.data.namespace = msg.namespace
		m.data.api = msg.api
		m.view = newContainersModel(m.size, m.data.namespace, m.msgCh)

		return m, m.view.Init()
	case containerLogsViewMsg:
		m.data.container = msg.container
		m.view = newContainerLogsModel(m.size, m.data.container, m.msgCh)

		return m, m.view.Init()
	case cronJobsViewMsg:
		m.data.namespace = msg.namespace
		m.data.api = msg.api
		m.view = newCronJobsModel(m.size, m.data.namespace, m.msgCh)

		return m, m.view.Init()
	case cronJobJobsViewMsg:
		m.data.cronJob = msg.cronJob
		m.view = newCronJobJobsModel(m.size, m.data.cronJob, m.msgCh)

		return m, m.view.Init()
	case cronJobContainersViewMsg:
		m.data.cronJobJob = msg.job
		m.view = newCronJobContainersModel(
			m.size,
			m.data.cronJob,
			m.data.cronJobJob,
			m.msgCh,
		)

		return m, m.view.Init()
	case cronJobLogsViewMsg:
		m.data.cronJobContainer = msg.container
		m.view = newCronJobLogsModel(
			m.size,
			m.data.cronJob,
			m.data.cronJobJob,
			m.data.cronJobContainer,
			m.msgCh,
		)
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
