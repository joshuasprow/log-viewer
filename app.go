package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
		key:   namespacesKey,
		view:  newNamespacesModel(size, msgCh),
	}
}

func (m appModel) Init() tea.Cmd {
	return m.view.Init()
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
		return newLogsModel(m.size, m.data.container, m.msgCh), nil
	case cronJobsKey:
		return newCronJobsModel(m.size, m.data.namespace, m.msgCh), nil
	case cronJobJobsKey:
		return newJobsModel(m.size, m.data.cronJob.Jobs, m.msgCh), nil
	case cronJobLogsKey:
		return newLogsModel(m.size, m.data.cronJobContainer, m.msgCh), nil
	default:
		return nil, fmt.Errorf("unknown view key: %s", key)
	}
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.Width = msg.Width
		m.size.Height = msg.Height - 1 // todo: fixes list title disappearing
	case viewMsg:
		var err error

		m.key = msg.key

		m.view, err = m.getViewByKey(m.key)
		if err != nil {
			panic(err)
		}

		m.data, err = updateViewData(m.data, msg)
		if err != nil {
			panic(err)
		}

		log.Printf("app.update - key=%s data=%+v", msg.key, msg.data)

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
