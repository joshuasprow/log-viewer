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
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"github.com/joshuasprow/log-viewer/pkg"
	"k8s.io/client-go/kubernetes"
)

type viewKey string

const (
	mainView       viewKey = ""
	namespacesView viewKey = "namespaces"
	podsView       viewKey = "pods"
)

var (
	clientset   *kubernetes.Clientset
	views       = map[viewKey]tea.Model{}
	defaultSize = tea.WindowSizeMsg{Width: 80, Height: 10}
)

func main() {
	cfg, err := loadConfig()
	check("load config", err)

	clientset, err = pkg.NewClientset(cfg.kubeconfig)
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

type podLogData string

func (d podLogData) FilterValue() string {
	return string(d)
}

type podData struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	Container string   `json:"container"`
	Logs      []string `json:"logs"`
}

func (d podData) FilterValue() string {
	if d.Container == "" {
		return d.Name
	}
	return fmt.Sprintf("%s/%s", d.Name, d.Container)
}

type namespaceData struct {
	Name string    `json:"name"`
	Pods []podData `json:"pods"`
}

func (d namespaceData) FilterValue() string {
	return d.Name
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
	pods, err := pkg.GetPodsNext(ctx, clientset, "")
	if err != nil {
		return nil, fmt.Errorf("get pods: %w", err)
	}

	namespaces := map[string]namespaceData{}

	for _, pod := range pods {
		pp := []podData{}

		if len(pod.Containers) == 0 {
			pp = append(pp, podData{
				Namespace: pod.Namespace,
				Name:      pod.Name,
				Logs:      []string{},
			})
		} else {
			for _, c := range pod.Containers {
				pp = append(pp, podData{
					Namespace: pod.Namespace,
					Name:      pod.Name,
					Container: c,
					Logs:      []string{},
				})
			}
		}

		namespace, ok := namespaces[pod.Namespace]

		if ok {
			namespace.Pods = append(namespace.Pods, pp...)
		} else {
			namespace = namespaceData{Name: pod.Namespace, Pods: pp}
		}

		namespaces[pod.Namespace] = namespace
	}

	data := []namespaceData{}

	for _, namespace := range namespaces {
		data = append(data, namespace)
	}

	return data, nil
}

type mainModel struct {
	model   list.Model
	spinner spinner.Model
	loading bool
	data    modelDataMsg
	err     error
	view    viewKey
	size    tea.WindowSizeMsg
}

type listItemDelegate struct{}

func (d listItemDelegate) Height() int                             { return 1 }
func (d listItemDelegate) Spacing() int                            { return 0 }
func (d listItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d listItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(list.Item)
	if !ok {
		return
	}

	fn := pkg.ListItemStyles.Normal.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return pkg.ListItemStyles.Selected.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(i.FilterValue()))
}

func newMainModel() mainModel {
	m := list.New(
		[]list.Item{},
		listItemDelegate{},
		defaultSize.Width,
		defaultSize.Height,
	)

	m.SetFilteringEnabled(false)
	m.SetShowStatusBar(false)

	m.Styles.PaginationStyle = pkg.ListStyles.Pagination
	m.Styles.HelpStyle = pkg.ListStyles.Help
	m.Styles.Title = pkg.ListStyles.Title
	m.Styles.TitleBar = pkg.ListStyles.TitleBar

	m.Title = "?"

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return mainModel{
		model:   m,
		spinner: s,
		loading: true,
	}
}

func (mainModel) initData() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		data, err := loadModelData(ctx)
		if err != nil {
			return errMsg{err: fmt.Errorf("load model data: %w", err)}
		}

		return modelDataMsg{data}
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.initData())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.err
		m.loading = false

		return m, nil
	case modelDataMsg:
		m.data = msg
		items := []list.Item{}

		for _, ns := range m.data.namespaces {
			items = append(items, ns)
		}

		m.model.SetItems(items)

		m.loading = false
	case tea.WindowSizeMsg:
		m.model.SetWidth(msg.Width)
		m.model.SetHeight(msg.Height)
		m.size = msg
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return mainModel{}, tea.Quit
		case "esc":
			if m.loading {
				return m, nil
			}

			switch m.view {
			case "":
			case namespacesView:
				m.view = ""
			case podsView:
				m.view = namespacesView
			}

			return m, nil
		case "enter":
			if m.loading {
				return m, nil
			}

			switch m.view {
			case "":
				selected := m.model.SelectedItem()

				namespace, ok := selected.(namespaceData)
				if !ok {
					m.err = fmt.Errorf("failed to cast %T as namespaceData", selected)
					return m, nil
				}

				view := newNamespaceModel(namespace, m.size)
				views[namespacesView] = view
				m.view = namespacesView

				return m, views[namespacesView].Init()
			case namespacesView:
				v, ok := views[namespacesView]
				if !ok {
					m.err = fmt.Errorf("failed to find namespace view")
				}

				n, ok := v.(namespaceModel)
				if !ok {
					m.err = fmt.Errorf("failed to cast %T as namespaceModel", v)
					return m, nil
				}

				selected := n.model.SelectedItem()

				pod, ok := selected.(podData)
				if !ok {
					m.err = fmt.Errorf("failed to cast %T as podData", selected)
					return m, nil
				}

				view := newPodModel(pod, m.size)
				views[podsView] = view
				m.view = podsView

				return m, views[podsView].Init()
			}
		}
	}

	var cmd tea.Cmd

	m.model, cmd = m.model.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	m.spinner, cmd = m.spinner.Update(msg)
	if cmd != nil {
		return m, cmd
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

	switch m.view {
	case "":
		if m.loading {
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().MarginRight(1).Render(m.spinner.View()),
				"loading...",
			)
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

type namespaceModel struct {
	data  namespaceData
	model list.Model
}

func newNamespaceModel(data namespaceData, size tea.WindowSizeMsg) namespaceModel {
	items := []list.Item{}

	for _, p := range data.Pods {
		items = append(items, p)
	}

	m := list.New(
		items,
		listItemDelegate{},
		size.Width,
		size.Height-1,
	)

	m.SetShowStatusBar(false)

	m.Title = fmt.Sprintf("%s:?", data.Name)

	m.Styles.PaginationStyle = pkg.ListStyles.Pagination
	m.Styles.HelpStyle = pkg.ListStyles.Help
	m.Styles.Title = pkg.ListStyles.Title
	m.Styles.TitleBar = pkg.ListStyles.TitleBar

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
	model   list.Model
	loading bool
	data    podData
}

func renderPodLogsTitle(data podData, loading bool) string {
	t := data.Namespace + ":"
	if data.Container == "" {
		t += data.Name
	} else {
		t += data.Name + "/" + data.Container
	}
	if loading {
		t += " - loading..."
	}
	return t
}

func newPodModel(data podData, size tea.WindowSizeMsg) podModel {
	items := []list.Item{}

	for _, l := range data.Logs {
		items = append(items, podLogData(l))
	}

	m := list.New(
		items,
		listItemDelegate{},
		size.Width,
		size.Height-1,
	)

	m.SetShowStatusBar(false)

	m.Styles.PaginationStyle = pkg.ListStyles.Pagination
	m.Styles.HelpStyle = pkg.ListStyles.Help
	m.Styles.Title = pkg.ListStyles.Title
	m.Styles.TitleBar = pkg.ListStyles.TitleBar

	loading := true

	m.Title = renderPodLogsTitle(data, loading)

	return podModel{
		model:   m,
		loading: loading,
		data:    data,
	}
}

type podLogsMsg []string

func (m podModel) Init() tea.Cmd {
	if m.data.Namespace == "" || m.data.Name == "" {
		return nil
	}

	return func() tea.Msg {
		ctx := context.Background()

		l, err := pkg.GetPodLogs(ctx, clientset, m.data.Namespace, m.data.Name, m.data.Container)
		if err != nil {
			return errMsg{fmt.Errorf("get pod logs: %w", err)}
		}

		return podLogsMsg(l)
	}
}

func (m podModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.model.SetWidth(msg.Width)
		m.model.SetHeight(msg.Height)
	case podLogsMsg:
		items := []list.Item{}

		for _, l := range msg {
			items = append(items, podLogData(l))
		}

		m.model.SetItems(items)
		m.loading = false
	}

	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m podModel) View() string {
	m.model.Title = renderPodLogsTitle(m.data, m.loading)

	return m.model.View()
}
