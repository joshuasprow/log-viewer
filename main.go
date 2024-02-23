package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/table"
	"github.com/kr/pretty"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
}

func run() {
	flags, err := parseFlags()
	check("parse flags", err)

	id := newResourceId(flags.namespace, flags.pod)

	clientset, err := newClientset(flags.kubeconfig)
	check("new clientset", err)

	ctx := context.Background()

	logCh, err := getPodLogsStream(ctx, clientset, id)
	check("get pd logs stream", err)

	rowCh := make(chan table.Row)

	go func() {
		defer close(rowCh)

		for log := range logCh {
			check("read log", log.err)

			pretty.Println(log.v)
		}
	}()

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

func newResourceId(namespace, pod string) resourceId {
	return resourceId{
		namespace: namespace,
		pod:       pod,
	}
}

func parseFlags() (flags, error) {
	return flags{
		kubeconfig: os.Getenv("KUBECONFIG"),
		namespace:  os.Getenv("NAMESPACE"),
		pod:        os.Getenv("POD"),
	}, nil
}

type resourceId struct {
	namespace string
	pod       string
}

type flags struct {
	kubeconfig string
	namespace  string
	pod        string
}

type logResult struct {
	v   map[string]any
	err error
}

func getPodLogsStream(ctx context.Context, clientset *kubernetes.Clientset, id resourceId) (<-chan logResult, error) {
	req := clientset.CoreV1().Pods(id.namespace).GetLogs(id.pod, &v1.PodLogOptions{})

	stream, err := req.Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("get stream: %w", err)
	}

	logCh := make(chan logResult)

	go func() {
		defer close(logCh)

		buf := make([]byte, 1024)

		for {
			n, err := stream.Read(buf)
			if err == io.EOF {
				return
			}
			if err != nil {
				logCh <- logResult{err: fmt.Errorf("read log: %w", err)}
				return
			}

			logCh <- logResult{v: map[string]any{"log": string(buf[:n])}}
		}
	}()

	return logCh, nil
}
