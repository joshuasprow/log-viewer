package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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

type podData struct {
	Name string   `json:"name"`
	Logs []string `json:"logs"`
}

type namespaceData struct {
	Name string    `json:"name"`
	Pods []podData `json:"pods"`
}

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

type modelDataMsg struct {
	namespaces []namespaceData
}

func loadModelData(ctx context.Context) ([]namespaceData, error) {
	pods, err := k8s.GetPodsNext(ctx, clientset, "")
	if err != nil {
		return nil, fmt.Errorf("get pods: %w", err)
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

	data := []namespaceData{}

	for _, namespace := range namespaces {
		data = append(data, namespace)
	}

	return data, nil
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
	model   list.Model
	loading bool
	data    modelDataMsg
	err     error
	view    viewKey
}

func newMainModel() mainModel {
	m := list.New([]list.Item{}, listItemDelegate{}, 10, 10)
	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(false)
	m.SetShowTitle(false)

	m.Styles.PaginationStyle = cli.ListStyles.Pagination
	m.Styles.HelpStyle = cli.ListStyles.Help

	return mainModel{
		model:   m,
		loading: true,
	}
}

func (m mainModel) Init() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		data, err := loadModelData(ctx)
		if err != nil {
			return errMsg{err: fmt.Errorf("load model data: %w", err)}
		}

		return modelDataMsg{data}
	}
}

func findNamespace(data modelDataMsg, namespace string) (namespaceData, error) {
	for _, ns := range data.namespaces {
		if ns.Name == namespace {
			return ns, nil
		}
	}

	return namespaceData{}, fmt.Errorf("namespace not found: %s", namespace)
}

func findPod(data modelDataMsg, namespace string, pod string) (podData, error) {
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
	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, nil
	case modelDataMsg:
		m.data = msg
		items := []list.Item{}

		for _, ns := range m.data.namespaces {
			items = append(items, listItem(ns.Name))
		}

		m.model.SetItems(items)

		m.loading = false
	case tea.WindowSizeMsg:
		m.model.SetWidth(msg.Width)
		m.model.SetHeight(msg.Height)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return mainModel{}, tea.Quit

		case "enter":
			if m.loading {
				return m, nil
			}

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
	if m.err != nil {
		return m.err.Error()
	}

	switch m.view {
	case "":
		if m.loading {
			return spinner.New().View()
		}

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
