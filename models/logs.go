package models

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"k8s.io/client-go/kubernetes"
)

type logListItem string

func (i logListItem) FilterValue() string {
	return string(i)
}

type LogsModel struct {
	clientset *kubernetes.Clientset
	model     list.Model
	namespace string
	pod       string
	container string
}

func Logs(
	clientset *kubernetes.Clientset,
	size tea.WindowSizeMsg,
	namespace string,
	pod string,
	container string,
) *LogsModel {
	m := defaultListModel(size)
	m.SetFilteringEnabled(true)
	m.Title = "logs"

	return &LogsModel{
		clientset: clientset,
		model:     m,
		namespace: namespace,
		pod:       pod,
		container: container,
	}
}

func (m *LogsModel) initData() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		logs, err := k8s.GetPodLogs(ctx, m.clientset, m.namespace, m.pod, m.container)
		if err != nil {
			return ErrMsg{Err: fmt.Errorf("load model data: %w", err)}
		}

		items := []list.Item{}

		for _, l := range logs {
			items = append(items, logListItem(l))
		}

		return LogsMsg(logs)
	}
}

func (m *LogsModel) Init() tea.Cmd {
	return tea.Batch(m.model.StartSpinner(), m.initData())
}

func (m *LogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrMsg:
		m.model.StopSpinner()
		return m, nil // todo: return error Cmd ?
	case LogsMsg:
		items := []list.Item{}

		for _, i := range msg {
			items = append(items, logListItem(i))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	case tea.WindowSizeMsg:
		m.model.SetWidth(msg.Width)
		m.model.SetHeight(msg.Height - 2)
	}

	var cmd tea.Cmd

	m.model, cmd = m.model.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

func (m *LogsModel) View() string {
	return m.model.View()
}

func (m *LogsModel) Selected() list.Item {
	return m.model.SelectedItem()
}
