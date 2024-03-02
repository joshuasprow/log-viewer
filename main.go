package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshuasprow/log-viewer/commands"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
	"github.com/joshuasprow/log-viewer/models"
	"github.com/joshuasprow/log-viewer/pkg"
	"k8s.io/client-go/kubernetes"
)

func main() {
	cfg, err := pkg.LoadConfig()
	check("load config", err)

	clientset, err := k8s.NewClientset(cfg.Kubeconfig)
	check("create k8s clientset", err)

	msgCh := make(chan tea.Msg)

	m := newMainModel(clientset, msgCh)

	prg := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for {
			msg := <-msgCh
			switch m := msg.(type) {
			case messages.Namespace:
				prg.Send(m)
			case messages.Container:
				prg.Send(m)
			}
		}
	}()

	_, err = prg.Run()

	check("run program", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}

type mainModel struct {
	clientset *kubernetes.Clientset
	msgCh     chan<- tea.Msg
	view      tea.Model
}

func newMainModel(
	clientset *kubernetes.Clientset,
	msgCh chan<- tea.Msg,
) mainModel {
	return mainModel{
		clientset: clientset,
		msgCh:     msgCh,
		view:      newNamespacesModel(clientset, msgCh),
	}
}

func (m mainModel) Init() tea.Cmd {
	return m.view.Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.Namespace:
		switch msg.View {
		case messages.NamespaceViewContainers:
			m.view = newContainersModel(m.clientset, msg.Name, m.msgCh)
			return m, m.view.Init()
		case messages.NamespaceViewCronJobs:
			m.view = newCronJobsModel(m.clientset, msg.Name, m.msgCh)
			return m, m.view.Init()
		default:
			panic(fmt.Errorf("unknown namespace view: %s", msg.View))
		}
	case messages.Container:
		m.view = newLogsModel(m.clientset, k8s.Container(msg), m.msgCh)
		return m, m.view.Init()
	}

	var cmd tea.Cmd

	if m.view != nil {
		m.view, cmd = m.view.Update(msg)
	}

	return m, cmd
}

func (m mainModel) View() string {
	var v string
	if m.view == nil {
		v = spinner.New().View()
	} else {
		v = m.view.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, v)
}

type namespaceListItem string

func (n namespaceListItem) FilterValue() string {
	return string(n)
}

type namespacesModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	msgCh     chan<- tea.Msg
}

func newNamespacesModel(
	clientset *kubernetes.Clientset,
	msgCh chan<- tea.Msg,
) namespacesModel {
	m := models.DefaultListModel()
	m.Title = "namespaces"

	return namespacesModel{
		clientset: clientset,
		model:     &m,
		msgCh:     msgCh,
	}
}

func (m namespacesModel) Init() tea.Cmd {
	return tea.Batch(
		m.model.StartSpinner(),
		commands.GetNamespaces(m.clientset),
	)
}

func (m namespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- messages.Namespace{Name: m.Selected()}
		}
	case messages.Namespaces:
		items := []list.Item{}

		for _, namespace := range msg {
			items = append(items, namespaceListItem(namespace))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm

	return m, cmd
}

func (m namespacesModel) View() string {
	return m.model.View()
}

func (m namespacesModel) Selected() string {
	return m.model.SelectedItem().(namespaceListItem).FilterValue()
}

type containerListItem k8s.Container

func (n containerListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s/%s", n.Namespace, n.Pod, n.Name)
}

type containersModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	namespace string
	msgCh     chan<- tea.Msg
}

func newContainersModel(
	clientset *kubernetes.Clientset,
	namespace string,
	msgCh chan<- tea.Msg,
) containersModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "containers"

	return containersModel{
		clientset: clientset,
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
}

func (m containersModel) Init() tea.Cmd {
	return tea.Batch(
		m.model.StartSpinner(),
		commands.GetCronJobs(m.clientset, m.namespace),
	)
}

func (m containersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- messages.Container(m.Selected())
		}
	case messages.Containers:
		items := []list.Item{}

		for _, c := range msg {
			items = append(items, containerListItem(c))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm
	return m, cmd
}

func (m containersModel) View() string {
	return m.model.View()
}

func (m containersModel) Selected() containerListItem {
	return m.model.SelectedItem().(containerListItem)
}

type cronJobListItem k8s.Container

func (n cronJobListItem) FilterValue() string {
	return fmt.Sprintf("%s/%s/%s", n.Namespace, n.Pod, n.Name)
}

type cronJobsModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	namespace string
	msgCh     chan<- tea.Msg
}

func newCronJobsModel(
	clientset *kubernetes.Clientset,
	namespace string,
	msgCh chan<- tea.Msg,
) cronJobsModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "cronJobs"

	return cronJobsModel{
		clientset: clientset,
		model:     &m,
		namespace: namespace,
		msgCh:     msgCh,
	}
}

func (m cronJobsModel) Init() tea.Cmd {
	return tea.Batch(
		m.model.StartSpinner(),
		commands.GetContainers(m.clientset, m.namespace),
	)
}

func (m cronJobsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.msgCh <- messages.Container(m.Selected())
		}
	case messages.Containers:
		items := []list.Item{}

		for _, c := range msg {
			items = append(items, cronJobListItem(c))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm
	return m, cmd
}

func (m cronJobsModel) View() string {
	return m.model.View()
}

func (m cronJobsModel) Selected() cronJobListItem {
	return m.model.SelectedItem().(cronJobListItem)
}

type logListItem string

func (i logListItem) FilterValue() string {
	return string(i)
}

type logsModel struct {
	clientset *kubernetes.Clientset
	model     *list.Model
	msgCh     chan<- tea.Msg

	namespace string
	pod       string
	container string
}

func newLogsModel(
	clientset *kubernetes.Clientset,
	container k8s.Container,
	msgCh chan<- tea.Msg,
) logsModel {
	m := models.DefaultListModel()
	m.SetFilteringEnabled(true)
	m.Title = "logs"

	return logsModel{
		clientset: clientset,
		model:     &m,
		msgCh:     msgCh,

		namespace: container.Namespace,
		pod:       container.Pod,
		container: container.Name,
	}
}

func (m logsModel) Init() tea.Cmd {
	return tea.Batch(
		m.model.StartSpinner(),
		commands.GetLogs(m.clientset, m.namespace, m.pod, m.container),
	)
}

func (m logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.Logs:
		items := []list.Item{}

		for _, i := range msg {
			items = append(items, logListItem(i))
		}

		m.model.SetItems(items)
		m.model.StopSpinner()
	}

	lm, cmd := m.model.Update(msg)
	m.model = &lm

	return m, cmd
}

func (m logsModel) View() string {
	return m.model.View()
}

func (m logsModel) Selected() list.Item {
	return m.model.SelectedItem()
}
