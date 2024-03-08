package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
	"github.com/joshuasprow/log-viewer/pkg"
	"github.com/joshuasprow/log-viewer/tui"
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
		models.Main(msgCh),
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
		case tui.NamespacesViewMsg:
			namespaces, err := k8s.GetNamespaces(ctx, clientset)
			check("get namespaces", err)

			prg.Send(tui.WrapNamespaces(namespaces))
		case tui.ApisViewMsg:
			prg.Send(tui.GetApis())
		case tui.ContainersViewMsg:
			containers, err := k8s.GetContainers(ctx, clientset, msg.Namespace, "")
			check("get containers", err)

			prg.Send(tui.WrapContainers(containers))
		case tui.ContainerLogsViewMsg:
			logs, err := k8s.GetPodLogs(ctx, clientset, msg.Container.Namespace, msg.Container.Pod, msg.Container.Name)
			check("get pod logs", err)

			prg.Send(tui.WrapLogs(logs))
		case tui.CronJobsViewMsg:
			cronJobs, err := k8s.GetCronJobs(ctx, clientset, msg.Namespace)
			check("get cron jobs", err)

			prg.Send(tui.WrapCronJobs(cronJobs))
		case tui.CronJobJobsViewMsg:
			jobs, err := k8s.GetJobs(ctx, clientset, msg.CronJob.Namespace, msg.CronJob.UID)
			check("get jobs", err)

			prg.Send(tui.WrapJobs(jobs))
		case tui.CronJobContainersViewMsg:
			labelSelector := fmt.Sprintf("job-name=%s", msg.Job.Name)

			containers, err := k8s.GetContainers(ctx, clientset, msg.Job.Namespace, labelSelector)
			check("get job containers", err)

			prg.Send(tui.WrapContainers(containers))
		case tui.CronJobLogsViewMsg:
			logs, err := k8s.GetPodLogs(ctx, clientset, msg.Container.Namespace, msg.Container.Pod, msg.Container.Name)
			check("get pod logs", err)

			prg.Send(tui.WrapLogs(logs))
		default:
			check("unknown message", fmt.Errorf("type %T", msg))
		}
	}
}
