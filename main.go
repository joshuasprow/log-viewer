package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/joshuasprow/log-viewer/k8s"
)

func main() {
	flags, err := parseFlags()
	check("parse flags", err)

	clientset, err := k8s.NewClientset(flags.kubeconfig)
	check("new clientset", err)

	p := tea.NewProgram(newModel())

	go func() {
		ctx := context.Background()

		namespaces, err := k8s.GetNamespaces(ctx, clientset)
		check("get namespaces", err)

		p.Send(namespacesMessage(namespaces))

		if len(namespaces) == 0 {
			return
		}
		pods, err := k8s.GetPods(ctx, clientset, namespaces[0])
		check("get pods", err)

		p.Send(podsMessage(pods))
	}()

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

type namespacesMessage []string
type podsMessage []string

type model struct {
	model table.Model
}

func newModel() tea.Model {
	t := table.New(
		table.WithFocused(true),
		table.WithKeyMap(table.DefaultKeyMap()),
	)

	t.SetColumns([]table.Column{
		{Title: "kind", Width: 12},
		{Title: "value", Width: 20},
	})

	t.SetStyles(table.DefaultStyles())

	return model{model: t}
}

func (m model) Init() tea.Cmd { return nil }

func toRows(kind string, values []string) []table.Row {
	rows := []table.Row{}

	for _, value := range values {
		rows = append(rows, table.Row{kind, value})
	}

	return rows
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case namespacesMessage:
		m.model.SetRows(append(m.model.Rows(), toRows("namespace", msg)...))
	case podsMessage:
		m.model.SetRows(append(m.model.Rows(), toRows("pod", msg)...))
	}

	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.model.View()
}
