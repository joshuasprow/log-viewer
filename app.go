package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func renderTitle(path ...string) string {
	var title string
	for i, p := range path {
		if i > 0 {
			title += " > "
		}
		color := "#FF00FF"
		if i == len(path)-1 {
			color = "#00FF00"
		}
		title += lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Render(p)
	}
	return title
}

type appModel struct {
	msgCh chan<- tea.Msg
	size  tea.WindowSizeMsg
	key   viewKey
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

func (m appModel) getViewByKey(key viewKey) (tea.Model, error) {
	switch key {
	case namespacesKey:
		return newNamespacesModel(m.size, m.msgCh), nil
	case apisKey:
		return newApisViewModel(m.size, m.data.namespace, m.msgCh), nil
	case containersKey:
		return newContainersModel(m.size, m.data.namespace, m.msgCh), nil
	case containerLogsKey:
		return newContainerLogsModel(m.size, m.data.container, m.msgCh), nil
	case cronJobsKey:
		return newCronJobsModel(m.size, m.data.namespace, m.msgCh), nil
	case cronJobJobsKey:
		return newCronJobJobsModel(
			m.size,
			m.data.cronJob,
			m.msgCh,
		), nil
	case cronJobContainersKey:
		return newCronJobContainersModel(m.size, m.data.cronJob, m.msgCh), nil
	case cronJobLogsKey:
		return newCronJobLogsModel(
			m.size,
			m.data.cronJobJob,
			m.msgCh,
		), nil
	default:
		return nil, fmt.Errorf("unknown view key: %s", key)
	}
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.Width = msg.Width
		m.size.Height = msg.Height - 1 // todo: fixes list title disappearing
	case namespacesViewMsg:
		m.key = namespacesKey
		m.view = newNamespacesModel(m.size, m.msgCh)

		return m, m.view.Init()
	case apisViewMsg:
		m.key = apisKey
		m.data.namespace = msg.namespace
		m.view = newApisViewModel(m.size, m.data.namespace, m.msgCh)

		return m, m.view.Init()
	case containersViewMsg:
		m.key = containersKey
		m.data.namespace = msg.namespace
		m.view = newContainersModel(m.size, m.data.namespace, m.msgCh)

		return m, m.view.Init()
	case containerLogsViewMsg:
		m.key = containerLogsKey
		m.data.container = msg.container
		m.view = newContainerLogsModel(m.size, m.data.container, m.msgCh)

		return m, m.view.Init()
	case cronJobsViewMsg:
		m.key = cronJobsKey
		m.data.namespace = msg.namespace
		m.view = newCronJobsModel(m.size, m.data.namespace, m.msgCh)

		return m, m.view.Init()
	case cronJobJobsViewMsg:
		m.key = cronJobJobsKey
		m.data.cronJob = msg.cronJob
		m.view = newCronJobJobsModel(m.size, m.data.cronJob, m.msgCh)

		return m, m.view.Init()
	case cronJobContainersViewMsg:
		m.key = cronJobContainersKey
		m.data.cronJobJob = msg.job
		m.view = newCronJobContainersModel(m.size, m.data.cronJob, m.msgCh)

		return m, m.view.Init()
	case cronJobLogsViewMsg:
		m.key = cronJobLogsKey
		m.data.cronJobContainer = msg.container
		m.view = newCronJobLogsModel(m.size, m.data.cronJobJob, m.msgCh)
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
