package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"k8s.io/client-go/kubernetes"
)

var views = map[ViewKey]tea.Model{}

type MainModel struct {
	clientset *kubernetes.Clientset
	view      ViewKey
	size      tea.WindowSizeMsg
	err       error
}

func Main(clientset *kubernetes.Clientset) MainModel {
	return MainModel{
		clientset: clientset,
		view:      NamespacesView,
		size:      defaultSize,
	}
}

func (m MainModel) Init() tea.Cmd {
	views[NamespacesView] = Namespaces(m.clientset, defaultSize)
	return views[NamespacesView].Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrMsg:
		m.err = msg.Err
		return m, nil
	case tea.WindowSizeMsg:
		m.size = msg
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return MainModel{}, tea.Quit
		case "esc":
			switch m.view {
			case ContainersView:
				m.view = NamespacesView
			case LogsView:
				m.view = ContainersView
			}
		case "enter":
			switch m.view {
			case NamespacesView:
				v, ok := views[NamespacesView]
				if !ok {
					m.err = fmt.Errorf("failed to find namespace view")
				}

				n, ok := v.(*NamespacesModel)
				if !ok {
					m.err = fmt.Errorf("failed to cast %T as *NamespacesModel", v)
					return m, nil
				}

				namespace := n.Selected()

				view := Containers(m.clientset, m.size, namespace)
				views[ContainersView] = view
				m.view = ContainersView

				return m, views[ContainersView].Init()

			case ContainersView:
				v, ok := views[ContainersView]
				if !ok {
					m.err = fmt.Errorf("failed to find containers view")
				}

				n, ok := v.(*ContainersModel)
				if !ok {
					m.err = fmt.Errorf("failed to cast %T as *ContainersModel", v)
					return m, nil
				}

				container := n.Selected()

				view := Logs(m.clientset, m.size, container.Namespace, container.Pod, container.Container)
				views[LogsView] = view
				m.view = LogsView

				return m, views[LogsView].Init()
			}
		}
	}

	var cmd tea.Cmd
	views[m.view], cmd = views[m.view].Update(msg)
	return m, cmd
}

func (m MainModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	view, ok := views[m.view]
	if !ok {
		return fmt.Sprintf("view not found: %s", m.view)
	}

	return view.View()
}
