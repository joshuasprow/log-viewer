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

var views = map[viewKey]tea.Model{
	namespacesView: namespacesModel{},
	podsView:       podsModel{},
}

type model struct {
	data modelData
	err  error
	view viewKey
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return model{}, tea.Quit

		case "enter":
			switch m.view {
			case "":
				if len(m.data.namespaces) == 0 {
					panic("no namespaces")
				}

				n := m.data.namespaces[0]

				namespace, err := findNamespace(m.data, n.Name)
				if err != nil {
					panic(fmt.Errorf("find namespace: %w", err))
				}

				view := newNamespacesModel(namespace)
				views[namespacesView] = view
				m.view = namespacesView

			case namespacesView:
				var nview namespacesModel
				v, ok := views[m.view]
				if !ok {
					panic(fmt.Errorf("view not found: %s", m.view))
				}

				nview, ok = v.(namespacesModel)
				if !ok {
					panic(fmt.Errorf("failed to cast %T as namespacesModel", v))
				}

				pod, err := findPod(m.data, nview.Name, nview.selected)
				if err != nil {
					panic(fmt.Errorf("find pod: %w", err))
				}

				pview := newPodsModel(pod)
				views[podsView] = pview
				m.view = podsView
			}
		}
	}

	for k, v := range views {
		v, _ := v.Update(msg)
		views[k] = v
	}

	return m, nil
}

func (m model) View() string {
	switch m.view {
	case "":
		return "hello"
	case namespacesView:
		view, ok := views[m.view]
		if !ok {
			return fmt.Sprintf("view not found: %s", m.view)
		}

		return view.View()
	case podsView:
		view, ok := views[m.view]
		if !ok {
			return fmt.Sprintf("view not found: %s", m.view)
		}

		return view.View()
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
			if len(m.Pods) == 0 {
				return m, nil
			}

			m.selected = m.Pods[0].Name

			return m, nil
		}
	}
	return m, nil
}

func (m namespacesModel) View() string {
	return m.Name + ":" + m.selected
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
			if len(m.Logs) == 0 {
				return m, nil
			}

			m.selected = m.Logs[0]

			return m, nil
		}
	}
	return m, nil
}

func (m podsModel) View() string {
	return m.Name + ":" + m.selected
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
