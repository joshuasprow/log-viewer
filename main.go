package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/joshuasprow/log-viewer/cli"
	"github.com/joshuasprow/log-viewer/k8s"
	"k8s.io/client-go/kubernetes"
)

var clientset *kubernetes.Clientset

func main() {
	cfg, err := loadConfig()
	check("load config", err)

	clientset, err = k8s.NewClientset(cfg.kubeconfig)
	check("create k8s clientset", err)

	ctx := context.Background()

	data, err := loadModelData(ctx)
	check("load model data", err)

	prg := tea.NewProgram(newMainModel(data))

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

func loadModelData(ctx context.Context) (modelData, error) {
	pods, err := k8s.GetPodsNext(ctx, clientset, "")
	if err != nil {
		return modelData{}, fmt.Errorf("get pods: %w", err)
	}

	namespaces := map[string]namespaceData{}

	for _, pod := range pods {
		p := podData{Name: pod.Name, Logs: []string{}}

		namespace, ok := namespaces[pod.Namespace]
		if ok {
			namespace.Pods = append(namespace.Pods, p)
		} else {
			namespace = namespaceData{
				Name: pod.Namespace,
				Pods: []podData{p},
			}
		}

		namespaces[pod.Namespace] = namespace
	}

	data := modelData{}

	for _, namespace := range namespaces {
		data.namespaces = append(data.namespaces, namespace)
	}

	return data, nil
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
	namespacesView: newNamespaceModel(namespaceData{}),
	podsView:       newPodModel(podData{}),
}

type mainModel struct {
	data  modelData
	model list.Model
	err   error
	view  viewKey
}

func newMainModel(data modelData) mainModel {
	items := []list.Item{}

	for _, p := range data.namespaces {
		items = append(items, listItem(p.Name))
	}

	m := list.New(items, listItemDelegate{}, 10, 10)
	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(false)
	m.SetShowTitle(false)

	m.Styles.PaginationStyle = cli.ListStyles.Pagination
	m.Styles.HelpStyle = cli.ListStyles.Help

	return mainModel{
		data:  data,
		model: m,
	}
}

func (m mainModel) Init() tea.Cmd { return nil }

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

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetWidth(msg.Width)
		m.model.SetHeight(msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return mainModel{}, tea.Quit

		case "enter":
			switch m.view {
			case "":
				if len(m.data.namespaces) == 0 {
					panic("no namespaces")
				}

				n := m.model.SelectedItem().FilterValue()

				namespace, err := findNamespace(m.data, n)
				if err != nil {
					panic(fmt.Errorf("find namespace: %w", err))
				}

				view := newNamespaceModel(namespace)
				views[namespacesView] = view
				m.view = namespacesView

			case namespacesView:
				var nview namespaceModel
				v, ok := views[m.view]
				if !ok {
					panic(fmt.Errorf("view not found: %s", m.view))
				}

				nview, ok = v.(namespaceModel)
				if !ok {
					panic(fmt.Errorf("failed to cast %T as namespacesModel", v))
				}

				selected := nview.model.SelectedItem().FilterValue()

				pod, err := findPod(m.data, nview.data.Name, selected)
				if err != nil {
					panic(fmt.Errorf("find pod: %w", err))
				}

				pview := newPodModel(pod)
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

func (m mainModel) View() string {
	switch m.view {
	case "":
		return m.model.View()
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

type listItem string

func (i listItem) FilterValue() string { return string(i) }

type listItemDelegate struct{}

func (d listItemDelegate) Height() int                             { return 1 }
func (d listItemDelegate) Spacing() int                            { return 0 }
func (d listItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(listItem)
	if !ok {
		return
	}

	fn := cli.ListItemStyles.Normal.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return cli.ListItemStyles.Selected.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(i.FilterValue()))
}

type namespaceModel struct {
	data  namespaceData
	model list.Model
}

func newNamespaceModel(data namespaceData) namespaceModel {
	items := []list.Item{}

	for _, p := range data.Pods {
		items = append(items, listItem(p.Name))
	}

	m := list.New(items, listItemDelegate{}, 10, 10)
	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(false)
	m.SetShowTitle(false)

	m.Styles.PaginationStyle = cli.ListStyles.Pagination
	m.Styles.HelpStyle = cli.ListStyles.Help

	return namespaceModel{
		data:  data,
		model: m,
	}
}

func (m namespaceModel) Init() tea.Cmd { return nil }

func (m namespaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetWidth(msg.Width)
		m.model.SetHeight(msg.Height)
	}

	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m namespaceModel) View() string {
	return m.model.View()
}

type podModel struct {
	data  podData
	model list.Model
}

func newPodModel(data podData) podModel {
	items := []list.Item{}

	for _, l := range data.Logs {
		items = append(items, listItem(l))
	}

	m := list.New(items, listItemDelegate{}, 10, 10)
	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(false)
	m.SetShowTitle(false)

	m.Styles.PaginationStyle = cli.ListStyles.Pagination
	m.Styles.HelpStyle = cli.ListStyles.Help

	return podModel{
		data:  data,
		model: m,
	}
}

func (m podModel) Init() tea.Cmd { return nil }

func (m podModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetWidth(msg.Width)
		m.model.SetHeight(msg.Height)
	}

	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m podModel) View() string {
	return m.model.View()
}
