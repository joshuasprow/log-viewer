package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	flags, err := parseFlags()
	check("parse flags", err)

	id := newResourceId(flags.namespace, flags.pod, flags.container)

	clientset, err := newClientset(flags.kubeconfig)
	check("new clientset", err)

	ctx := context.Background()

	logCh, err := getPodLogsStream(ctx, clientset, id)
	check("get pd logs stream", err)

	for log := range logCh {
		check("read log", log.err)

		pretty.Println(log.v)
	}

	// m := cli.NewModel(rowCh)

	// _, err = tea.NewProgram(m).Run()
	// check("run program", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}

func newClientset(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func newResourceId(namespace, pod, container string) resourceId {
	return resourceId{
		namespace: namespace,
		pod:       pod,
		container: container,
	}
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

	namespace := os.Getenv("NAMESPACE")
	pod := os.Getenv("POD")
	container := os.Getenv("CONTAINER")

	return flags{
		kubeconfig: kubeconfig,
		namespace:  namespace,
		pod:        pod,
		container:  container,
	}, nil
}

type resourceId struct {
	namespace string
	pod       string
	container string
}

type flags struct {
	kubeconfig string
	namespace  string
	pod        string
	container  string
}

type logResult struct {
	v   map[string]any
	err error
}

func ptr[V any](v V) *V {
	return &v
}

func getPodLogs(ctx context.Context, clientset *kubernetes.Clientset, id resourceId) (string, error) {
	body, err := clientset.
		CoreV1().
		Pods(id.namespace).
		GetLogs(
			id.pod,
			&v1.PodLogOptions{
				Container: id.container,
			},
		).
		Do(ctx).
		Raw()

	return string(body), err
}

func getPodLogsStream(ctx context.Context, clientset *kubernetes.Clientset, id resourceId) (<-chan logResult, error) {
	req := clientset.
		CoreV1().
		Pods(id.namespace).
		GetLogs(
			id.pod,
			&v1.PodLogOptions{
				Container: id.container,
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

			b := scanner.Bytes()
			v := map[string]any{}

			if err := json.Unmarshal(b, &v); err != nil {
				logCh <- logResult{err: fmt.Errorf("unmarshal log %q: %w", string(b), err)}
				return
			}

			logCh <- logResult{v: v}
		}
	}()

	return logCh, nil
}
