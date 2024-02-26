package models

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/pkg"
	"k8s.io/client-go/kubernetes"
)

type namespaceListItem string

func (n namespaceListItem) FilterValue() string {
	return string(n)
}

type NamespacesModel struct {
	clientset *kubernetes.Clientset
	model     list.Model
}

func Namespaces(clientset *kubernetes.Clientset, size tea.WindowSizeMsg) tea.Model {
	m := list.New(
		[]list.Item{},
		listItemDelegate{},
		size.Width,
		size.Height-2,
	)

	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(false)

	m.Styles.PaginationStyle = listStyles.Pagination
	m.Styles.HelpStyle = listStyles.Help
	m.Styles.Title = listStyles.Title
	m.Styles.TitleBar = listStyles.TitleBar

	m.Title = "namespaces"

	return &NamespacesModel{
		clientset: clientset,
		model:     m,
	}
}

func (m *NamespacesModel) initData() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		namespaces, err := pkg.GetNamespaces(ctx, m.clientset)
		if err != nil {
			return ErrMsg{Err: fmt.Errorf("load model data: %w", err)}
		}

		return NamespacesMsg(namespaces)
	}
}

func (m *NamespacesModel) Init() tea.Cmd {
	return tea.Batch(m.model.StartSpinner(), m.initData())
}

func (m *NamespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrMsg:
		m.model.StopSpinner()
		return m, nil // todo: return error Cmd ?
	case NamespacesMsg:
		items := []list.Item{}

		for _, namespace := range msg {
			items = append(items, namespaceListItem(namespace))
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

func (m *NamespacesModel) View() string {
	return m.model.View()
}

func (m *NamespacesModel) Selected() string {
	return m.model.SelectedItem().(namespaceListItem).FilterValue()
}
