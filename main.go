package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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
	msgCh := make(chan tea.Msg)

	logFile, err := tea.LogToFile("tmp/debug.log", "")
	check("log to file", err)
	defer logFile.Close()

	log.SetOutput(logFile)

	prg := tea.NewProgram(
		newAppModel(msgCh),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	)

	go handleMessages(ctx, clientset, prg, msgCh)

	_, err = prg.Run()
	check("run program", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}

func handleMessages(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	prg *tea.Program,
	msgCh <-chan tea.Msg,
) {
	for msg := range msgCh {
		prg.Send(msg)

		switch msg := msg.(type) {
		case viewMsg,
			messages.CronJob:
			// ...
		case messages.Init:
			namespaces, err := k8s.GetNamespaces(ctx, clientset)
			check("get namespaces", err)

			prg.Send(messages.Namespaces(namespaces))
		case messages.Namespace:
			switch msg.Api {
			case "":

			case messages.ContainersApi:
				containers, err := k8s.GetContainers(ctx, clientset, msg.Name, "")
				check("get containers", err)

				prg.Send(messages.Containers(containers))
			case messages.CronJobsApi:
				cronJobs, err := k8s.GetCronJobs(ctx, clientset, msg.Name)
				check("get cron jobs", err)

				prg.Send(messages.CronJobs(cronJobs))
			}
		case messages.Container:
			logs, err := k8s.GetPodLogs(ctx, clientset, msg.Namespace, msg.Pod, msg.Name)
			check("get pod logs", err)

			prg.Send(messages.Logs(logs))
		case messages.JobContainer:
			logs, err := k8s.GetPodLogs(ctx, clientset, msg.Namespace, msg.Pod, msg.Name)
			check("get pod logs", err)

			prg.Send(messages.Logs(logs))
		case messages.Job:
			labelSelector := fmt.Sprintf("job-name=%s", msg.Name)

			containers, err := k8s.GetContainers(ctx, clientset, msg.Namespace, labelSelector)
			check("get containers", err)

			prg.Send(messages.JobContainers(containers))
		default:
			check("unknown message", fmt.Errorf("type %T", msg))
		}
	}
}
