package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
	"k8s.io/client-go/kubernetes"
)

var (
	clientset   *kubernetes.Clientset
	defaultSize = tea.WindowSizeMsg{Width: 80, Height: 10}
	views       = map[models.View]tea.Model{}
)

func main() {
	cfg, err := loadConfig()
	check("load config", err)

	clientset, err = k8s.NewClientset(cfg.kubeconfig)
	check("create k8s clientset", err)

	prg := tea.NewProgram(newMainModel())

	_, err = prg.Run()
	check("run program", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}

type config struct {
	kubeconfig string
}

func loadConfig() (config, error) {
	godotenv.Load()

	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return config{}, fmt.Errorf("get user home dir: %w", err)
		}

		kubeconfig = filepath.Join(homedir, ".kube", "config")
	}

	return config{kubeconfig: kubeconfig}, nil
}

type mainModel struct {
	err  error
	view models.View
	size tea.WindowSizeMsg
}

func newMainModel() mainModel {
	return mainModel{
		view: models.NamespacesView,
		size: defaultSize,
	}
}

func (m mainModel) Init() tea.Cmd {
	views[models.NamespacesView] = models.Namespaces(clientset, defaultSize)
	return views[models.NamespacesView].Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case models.ErrMsg:
		m.err = msg.Err
		return m, nil
	case tea.WindowSizeMsg:
		m.size = msg
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return mainModel{}, tea.Quit
		case "esc":
			switch m.view {
			case models.ContainersView:
				m.view = models.NamespacesView
			case models.LogsView:
				m.view = models.ContainersView
			}
			return m, nil
		case "enter":
			switch m.view {
			case models.NamespacesView:
				v, ok := views[models.NamespacesView]
				if !ok {
					m.err = fmt.Errorf("failed to find namespace view")
				}

				n, ok := v.(*models.NamespacesModel)
				if !ok {
					m.err = fmt.Errorf("failed to cast %T as *models.NamespacesModel", v)
					return m, nil
				}

				namespace := n.Selected()

				view := models.Containers(clientset, m.size, namespace)
				views[models.ContainersView] = view
				m.view = models.ContainersView

				return m, views[models.ContainersView].Init()

			case models.ContainersView:
				v, ok := views[models.ContainersView]
				if !ok {
					m.err = fmt.Errorf("failed to find containers view")
				}

				n, ok := v.(*models.ContainersModel)
				if !ok {
					m.err = fmt.Errorf("failed to cast %T as *models.ContainersModel", v)
					return m, nil
				}

				container := n.Selected()

				view := models.Logs(clientset, m.size, container.Namespace, container.Pod, container.Container)
				views[models.LogsView] = view
				m.view = models.LogsView

				return m, views[models.LogsView].Init()
			}
		}
	}

	for k, v := range views {
		v, cmd := v.Update(msg)
		if cmd != nil {
			return m, cmd
		}

		views[k] = v
	}

	return m, nil
}

func (m mainModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	view, ok := views[m.view]
	if !ok {
		return fmt.Sprintf("view not found: %s", m.view)
	}

	return view.View()
}
