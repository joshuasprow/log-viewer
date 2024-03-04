package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
	"github.com/joshuasprow/log-viewer/pkg"
)

func main() {
	cfg, err := pkg.LoadConfig()
	check("load config", err)

	clientset, err := k8s.NewClientset(cfg.Kubeconfig)
	check("create k8s clientset", err)

	msgCh := make(chan tea.Msg)
	ctx := context.Background()

	prg := tea.NewProgram(
		newAppModel(clientset, msgCh),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	)

	go func() {
		for {
			msg := <-msgCh
			switch m := msg.(type) {
			case messages.Init:
				namespaces, err := k8s.GetNamespaces(ctx, clientset)
				check("get namespaces", err)

				prg.Send(messages.Namespaces(namespaces))
			case messages.Namespace,
				messages.Container,
				messages.CronJob,
				messages.Job,
				messages.JobContainer:
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
