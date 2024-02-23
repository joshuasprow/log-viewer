package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/pkg"
	"k8s.io/client-go/kubernetes"
)

type viewKey string

const (
	namespaceKey viewKey = "namespace"
	podKey       viewKey = "pod"
	logsKey      viewKey = "logs"
)

type model struct {
	clientset      *kubernetes.Clientset
	view           viewKey
	namespaceModel tea.Model
	namespace      string
	podModel       tea.Model
	pod            string
	logsModel      tea.Model
	logsCh         chan pkg.LogEntry
	err            error
}

func NewModel(clientset *kubernetes.Clientset) tea.Model {
	return model{
		clientset:      clientset,
		view:           namespaceKey,
		namespaceModel: newNamespacesModel(clientset),
		podModel:       newListModel([]string{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func getSelectedListItem(m tea.Model) (string, error) {
	var lm listModel

	switch t := m.(type) {
	case namespacesModel:
		lm = t.model.(listModel)
	case listModel:
		lm = t
	default:
		return "", fmt.Errorf("unexpected model type: %T", m)
	}

	i, ok := lm.l.SelectedItem().(listItem)
	if !ok {
		return "", fmt.Errorf("expected listModel.SelectedItem to return an item, but got %T", lm.l.SelectedItem())
	}

	return string(i), nil
}

func (m model) updateViews(msg tea.Msg) (model, tea.Cmd) {
	var nCmd tea.Cmd
	m.namespaceModel, nCmd = m.namespaceModel.Update(msg)

	var pCmd tea.Cmd
	m.podModel, pCmd = m.podModel.Update(msg)

	var lCmd tea.Cmd
	if m.logsModel != nil {
		m.logsModel, lCmd = m.logsModel.Update(msg)
	}

	return m, tea.Batch(nCmd, pCmd, lCmd)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		panic(fmt.Errorf("model.err: %w", m.err))
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		panic(fmt.Errorf("model.Update: %w", msg))
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.view {
			case namespaceKey:
				namespace, err := getSelectedListItem(m.namespaceModel)
				if err != nil {
					m.err = fmt.Errorf("failed to get selected namespace: %w", err)
					return m, nil
				}

				m.namespace = namespace

				pods, ok := pkg.Pods[m.namespace]
				if !ok {
					m.err = fmt.Errorf("no pods found in namespace %s", m.namespace)
					return m, nil
				}

				m.podModel = newListModel(pods)
				m.view = podKey
			case podKey:
				pod, err := getSelectedListItem(m.podModel)
				if err != nil {
					m.err = fmt.Errorf("failed to get selected pod: %w", err)
					return m, nil
				}

				if m.namespace == "" {
					m.err = fmt.Errorf("no namespace selected")
					return m, nil
				}

				if pod == "" {
					m.err = fmt.Errorf("no pod selected")
					return m, nil
				}

				m.pod = pod
				m.view = logsKey
				m.logsModel = newLogsModel(m.clientset, k8s.NewResourceId(m.namespace, m.pod, ""))

				cmd = m.logsModel.Init()
			}
		case "esc":
			switch m.view {
			case podKey:
				m.pod = ""
				m.podModel = newListModel([]string{})

				m.view = namespaceKey
			case logsKey:
				// todo: clean up logs stream and model
				m.view = podKey
			}
		}
	}

	m, cmd = m.updateViews(msg)

	return m, cmd
}

func (m model) View() string {
	body := ""

	switch m.view {
	case namespaceKey:
		body = m.namespaceModel.View()
	case podKey:
		body = m.podModel.View()
	case logsKey:
		body = m.logsModel.View()
	default:
		m.err = fmt.Errorf("unknown view key: %q", m.view)

		return m.err.Error()
	}

	titlePart := func(k, v string) string {
		return lipgloss.NewStyle().PaddingLeft(2).Render(fmt.Sprintf("%s: %q", k, v))
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			titlePart("view", string(m.view)),
			titlePart("source", m.namespace),
			titlePart("resource", m.pod),
		),
		body,
	)
}
