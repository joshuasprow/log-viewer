package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type mainModel struct {
	msgCh chan<- tea.Msg
	size  tea.WindowSizeMsg
	view  tea.Model
	data  viewData
	err   error
}

func newMainModel(msgCh chan<- tea.Msg) mainModel {
	size := tea.WindowSizeMsg{Width: 80, Height: 24}

	return mainModel{
		msgCh: msgCh,
		size:  size,
	}
}

func (m mainModel) Init() tea.Cmd {
	return func() tea.Msg {
		m.msgCh <- namespacesViewMsg{}
		return nil
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		m.view = newErrorModel(m.size, m.err, m.msgCh)
		m.err = nil
		return m, nil
	}

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

	if m.view == nil {
		return m, nil
	}

	var cmd tea.Cmd
	m.view, cmd = m.view.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	if m.view == nil {
		return spinner.New().View()
	}
	return m.view.View()
}