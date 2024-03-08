package models

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/tui"
)

type mainModel struct {
	msgCh chan<- tea.Msg
	size  tea.WindowSizeMsg
	view  tea.Model
	data  tui.ViewData
	err   error
}

func Main(msgCh chan<- tea.Msg) mainModel {
	size := tea.WindowSizeMsg{Width: 80, Height: 24}

	return mainModel{
		msgCh: msgCh,
		size:  size,
	}
}

func (m mainModel) Init() tea.Cmd {
	return func() tea.Msg {
		m.msgCh <- tui.NamespacesViewMsg{}
		return nil
	}
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		m.view = Error(m.size, m.err, m.msgCh)
		m.err = nil
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.Width = msg.Width
		m.size.Height = msg.Height - 1 // todo: fixes list title disappearing
	case tui.NamespacesViewMsg:
		m.view = Namespaces(m.size, m.msgCh)
		return m, m.view.Init()
	case tui.ApisViewMsg:
		m.data.Namespace = msg.Namespace
		m.view = Apis(m.size, m.data.Namespace, m.msgCh)
		return m, m.view.Init()
	case tui.ContainersViewMsg:
		m.data.Namespace = msg.Namespace
		m.data.Api = msg.Api
		m.view = Containers(m.size, m.data.Namespace, m.msgCh)
		return m, m.view.Init()
	case tui.ContainerLogsViewMsg:
		m.data.Container = msg.Container
		m.view = ContainerLogs(m.size, m.data.Container, m.msgCh)
		return m, m.view.Init()
	case tui.CronJobsViewMsg:
		m.data.Namespace = msg.Namespace
		m.data.Api = msg.Api
		m.view = CronJobs(m.size, m.data.Namespace, m.msgCh)
		return m, m.view.Init()
	case tui.CronJobJobsViewMsg:
		m.data.CronJob = msg.CronJob
		m.view = CronJobJobs(m.size, m.data.CronJob, m.msgCh)
		return m, m.view.Init()
	case tui.CronJobContainersViewMsg:
		m.data.CronJobJob = msg.Job
		m.view = CronJobContainers(
			m.size,
			m.data.CronJob,
			m.data.CronJobJob,
			m.msgCh,
		)
		return m, m.view.Init()
	case tui.CronJobLogsViewMsg:
		m.data.CronJobContainer = msg.Container
		m.view = CronJobLogs(
			m.size,
			m.data.CronJob,
			m.data.CronJobJob,
			m.data.CronJobContainer,
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
