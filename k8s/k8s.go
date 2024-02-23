package k8s

import (
	"bufio"
	"context"
	"fmt"

	"github.com/joshuasprow/log-viewer/pkg"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClientset(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

type ResourceId struct {
	Namespace string
	Pod       string
	Container string
}

func NewResourceId(namespace, pod, container string) ResourceId {
	return ResourceId{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
	}
}

func getPodLogs(ctx context.Context, clientset *kubernetes.Clientset, id ResourceId) (string, error) {
	body, err := clientset.
		CoreV1().
		Pods(id.Namespace).
		GetLogs(
			id.Pod,
			&v1.PodLogOptions{
				Container: id.Container,
			},
		).
		Do(ctx).
		Raw()

	return string(body), err
}

func StreamPodLogs(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	id ResourceId,
	logsCh chan<- pkg.Result[pkg.LogEntry],
) {
	defer close(logsCh)

	type R = pkg.Result[pkg.LogEntry]

	req := clientset.
		CoreV1().
		Pods(id.Namespace).
		GetLogs(
			id.Pod,
			&v1.PodLogOptions{
				Container: id.Container,
				Follow:    true,
				TailLines: pkg.Ptr(int64(10)),
			},
		)

	stream, err := req.Stream(ctx)
	if err != nil {
		logsCh <- R{Err: fmt.Errorf("get stream: %w", err)}
		return
	}
	defer func() {
		if err := stream.Close(); err != nil {
			logsCh <- R{Err: fmt.Errorf("close stream: %w", err)}
		}
	}()

	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			logsCh <- R{Err: fmt.Errorf("scan error: %w", err)}
			return
		}

		data := scanner.Bytes()

		v, err := pkg.UnmarshalLogEntry(data)
		if err != nil {
			logsCh <- R{Err: fmt.Errorf("unmarshal log entry %q: %w", string(data), err)}
			return
		}

		logsCh <- R{V: v}
	}
}
