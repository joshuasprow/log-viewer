package k8s

import (
	"bufio"
	"context"
	"fmt"

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

func newResourceId(namespace, pod, container string) ResourceId {
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

func GetPodLogsStream(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	id ResourceId,
) (
	<-chan logResult,
	error,
) {
	req := clientset.
		CoreV1().
		Pods(id.Namespace).
		GetLogs(
			id.Pod,
			&v1.PodLogOptions{
				Container: id.Container,
				Follow:    true,
				TailLines: ptr(int64(10)),
			},
		)

	stream, err := req.Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("get stream: %w", err)
	}

	logCh := make(chan logResult)

	scanner := bufio.NewScanner(stream)

	go func() {
		defer close(logCh)

		defer func() {
			if err := stream.Close(); err != nil {
				logCh <- logResult{err: fmt.Errorf("close stream: %w", err)}
			}
		}()

		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				logCh <- logResult{err: fmt.Errorf("scan error: %w", err)}
				return
			}

			data := scanner.Bytes()

			logCh <- logResult{v: v}
		}
	}()

	return logCh, nil
}
