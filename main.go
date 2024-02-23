package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/joshuasprow/log-viewer/cli"
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

	rowCh := make(chan cli.TableRowItem)

	go func() {
		defer close(rowCh)

		for log := range logCh {
			check("read log", log.err)

			rowCh <- log.v
		}
	}()

	m := cli.NewModel(rowCh)

	_, err = tea.NewProgram(m).Run()
	check("run program", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
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

	return flags{kubeconfig: kubeconfig}, nil
}

type flags struct {
	kubeconfig string
}
