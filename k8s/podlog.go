package k8s

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/joshuasprow/log-viewer/pkg"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPodLogs(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	pod string,
	container string,
) (
	[]string,
	error,
) {
	data, err := clientset.
		CoreV1().
		Pods(namespace).
		GetLogs(
			pod,
			&v1.PodLogOptions{
				Container: container,
				TailLines: pkg.Ptr(int64(10)),
			},
		).
		Do(ctx).
		Raw()
	if err != nil {
		return nil, err
	}

	logs := []string{}

	for _, line := range strings.Split(string(data), "\n") {
		l := strings.TrimSpace(line)

		if l == "" {
			continue
		}

		logs = append(logs, l)
	}

	return logs, nil
}

// todo: stream logs when view is initialized
func StreamPodLogs(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	pod string,
	container string,
	logsCh chan<- pkg.Result[string],
) {
	defer close(logsCh)

	type R = pkg.Result[string]

	req := clientset.
		CoreV1().
		Pods(namespace).
		GetLogs(
			pod,
			&v1.PodLogOptions{
				Container: container,
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

		logsCh <- R{V: scanner.Text()}
	}
}
