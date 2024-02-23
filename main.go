package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/joshuasprow/log-viewer/cli"
	"github.com/joshuasprow/log-viewer/k8s"
)

func main() {
	flags, err := parseFlags()
	check("parse flags", err)

	clientset, err := k8s.NewClientset(flags.kubeconfig)
	check("new clientset", err)

	m := cli.NewModel(clientset)

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
