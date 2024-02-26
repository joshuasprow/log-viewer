package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func main() {
	d, err := readModelData()
	check("read model data", err)

	p := tea.NewProgram(newModel(d))

	_, err = p.Run()
	check("run program", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}

type flags struct {
	kubeconfig string
}

func parseFlags() (flags, error) {
	godotenv.Load()

	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return flags{}, fmt.Errorf("get user home dir: %w", err)
		}

		kubeconfig = filepath.Join(homedir, ".kube", "config")
	}

	return flags{kubeconfig: kubeconfig}, nil
}

type podData struct {
	Name string   `json:"name"`
	Logs []string `json:"logs"`
}

type namespaceData struct {
	Name string    `json:"name"`
	Pods []podData `json:"pods"`
}

type modelData struct {
	namespaces []namespaceData
}

func readModelData() (modelData, error) {
	d, err := os.ReadFile("tmp/data.json")
	if err != nil {
		return modelData{}, fmt.Errorf("read file: %w", err)
	}

	namespaces := []namespaceData{}

	err = json.Unmarshal(d, &namespaces)
	if err != nil {
		return modelData{}, fmt.Errorf("unmarshal data: %w", err)
	}

	return modelData{namespaces: namespaces}, nil
}

type viewKey string

const (
	mainView       viewKey = ""
	namespacesView viewKey = "namespaces"
	podsView       viewKey = "pods"
)

var views = map[string]childModel{
	string(namespacesView): {},
	string(podsView):       {},
}

type model struct {
	data modelData
	err  error
	view string
}

func newModel(data modelData) model {
	return model{data: data}
}

func (m model) Init() tea.Cmd { return nil }

func findNamespace(data modelData, namespace string) (namespaceData, error) {
	for _, ns := range data.namespaces {
		if ns.Name == namespace {
			return ns, nil
		}
	}

	return namespaceData{}, fmt.Errorf("namespace not found: %s", namespace)
}

func findPod(data modelData, namespace string, pod string) (podData, error) {
	ns, err := findNamespace(data, namespace)
	if err != nil {
		return podData{}, fmt.Errorf("find namespace: %w", err)
	}

	for _, p := range ns.Pods {
		if p.Name == pod {
			return p, nil
		}
	}

	return podData{}, fmt.Errorf("pod not found: %s", pod)
}

func setView(m model, viewKey string, viewName string) (model, error) {
	view, ok := views[viewKey]
	if !ok {
		return m, fmt.Errorf("view not found: %s", viewKey)
	}

	m.view = viewKey
	view.name = viewName
	views[viewKey] = view

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return model{}, tea.Quit

		case "enter":
			switch m.view {
			case "":
				n := m.data.namespaces[0]

				namespace, err := findNamespace(m.data, n.Name)
				if err != nil {
					panic(fmt.Errorf("find namespace: %w", err))
				}

				m, err = setView(m, string(namespacesView), namespace.Name)
				if err != nil {
					panic(fmt.Errorf("set view: %w", err))
				}
			case string(namespacesView):
				view, ok := views[m.view]
				if !ok {
					panic(fmt.Errorf("view not found: %s", m.view))
				}

				namespace, err := findNamespace(m.data, view.name)
				if err != nil {
					panic(fmt.Errorf("find namespace: %w", err))
				}

				pod, err := findPod(m.data, namespace.Name, namespace.Pods[0].Name)
				if err != nil {
					panic(fmt.Errorf("find pod: %w", err))
				}

				m, err = setView(m, string(podsView), pod.Name)
				if err != nil {
					panic(fmt.Errorf("set view: %w", err))
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.view {
	case "":
		return "hello"
	case string(namespacesView):
		view, ok := views[m.view]
		if !ok {
			return fmt.Sprintf("view not found: %s", m.view)
		}

		return view.name
	case string(podsView):
		view, ok := views[m.view]
		if !ok {
			return fmt.Sprintf("view not found: %s", m.view)
		}

		return view.name
	default:
		return fmt.Sprintf("view not found: %s", m.view)
	}
}

type namespacesModel struct {
	namespaceData
	selected string
}

func newNamespacesModel(data namespaceData) namespacesModel {
	return namespacesModel{namespaceData: data}
}

func (m namespacesModel) Init() tea.Cmd { return nil }

func (m namespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.selected = m.Pods[0].Name
			return m, nil
		}
	}
	return m, nil
}

func (m namespacesModel) View() string {
	return "hello"
}

type podsModel struct {
	podData
	selected string
}

func newPodsModel(data podData) podsModel {
	return podsModel{podData: data}
}

func (m podsModel) Init() tea.Cmd { return nil }

func (m podsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.selected = m.Logs[0]
			return m, nil
		}
	}
	return m, nil
}

func (m podsModel) View() string {
	return "hello"
}

type childModel struct {
	name  string
	names []string

	selected string
}

func newChildModel(name string, names []string) childModel {
	return childModel{}
}

func (m childModel) Init() tea.Cmd { return nil }

func (m childModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			return m, nil
		}
	}
	return m, nil
}

func (m childModel) View() string {
	return "hello"
}
