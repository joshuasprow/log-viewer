package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
	"github.com/joshuasprow/log-viewer/pkg"
	"k8s.io/client-go/kubernetes"
)

func main() {
	cfg, err := pkg.LoadConfig()
	check("load config", err)

	clientset, err := k8s.NewClientset(cfg.Kubeconfig)
	check("create k8s clientset", err)

	ctx := context.Background()

	cronJobs, err := k8s.GetCronJobs(ctx, clientset, "default")
	check("list jobs", err)

	for _, cronJob := range cronJobs {
		fmt.Printf("%s (%d)\n", cronJob.Name, len(cronJob.Jobs))
	}
	// msgCh := make(chan tea.Msg)

	// m := newMainModel(clientset, msgCh)

	// prg := tea.NewProgram(m, tea.WithAltScreen())

	// go func() {
	// 	for {
	// 		msg := <-msgCh
	// 		switch m := msg.(type) {
	// 		case messages.Namespace:
	// 			prg.Send(m)
	// 		case messages.Container:
	// 			prg.Send(m)
	// 		case messages.CronJob:
	// 			prg.Send(m)
	// 		}
	// 	}
	// }()

	// _, err = prg.Run()

	// check("run program", err)
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
		switch msg.Api {
		case "":
			m.view = newApisViewModel(m.clientset, msg.Name, m.msgCh)
			return m, m.view.Init()
		case messages.ContainersApi:
			m.view = newContainersModel(m.clientset, msg.Name, m.msgCh)
			return m, m.view.Init()
		case messages.CronJobsApi:
			m.view = newCronJobsModel(m.clientset, msg.Name, m.msgCh)
			return m, m.view.Init()
		default:
			panic(fmt.Errorf("unknown namespace view: %s", msg.Api))
		}
	case messages.Container:
		m.view = newLogsModel(m.clientset, k8s.Container(msg), m.msgCh)
		return m, m.view.Init()
	case messages.CronJob:
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
