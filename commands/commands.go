package commands

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/messages"
	"k8s.io/client-go/kubernetes"
)

func GetNamespaces(clientset *kubernetes.Clientset) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		namespaces, err := k8s.GetNamespaces(ctx, clientset)
		if err != nil {
			return messages.Error{
				Err: fmt.Errorf("get namespaces: %w", err),
			}
		}

		return messages.Namespaces(namespaces)
	}
}

func GetContainers(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		containers, err := k8s.GetContainers(ctx, clientset, namespace)
		if err != nil {
			return messages.Error{
				Err: fmt.Errorf("get containers: %w", err),
			}
		}

		return messages.Containers(containers)
	}
}

func GetCronJobs(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		cronJobs, err := k8s.GetCronJobs(ctx, clientset, namespace)
		if err != nil {
			return messages.Error{
				Err: fmt.Errorf("get cron jobs: %w", err),
			}
		}

		return messages.CronJobs(cronJobs)
	}
}

func GetLogs(
	clientset *kubernetes.Clientset,
	namespace,
	pod string,
	container string,
) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		logs, err := k8s.GetPodLogs(ctx, clientset, namespace, pod, container)
		if err != nil {
			return messages.Error{
				Err: fmt.Errorf("get pod logs: %w", err),
			}
		}

		return messages.Logs(logs)
	}
}
