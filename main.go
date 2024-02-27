package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/models"
)

var (
	defaultSize = tea.WindowSizeMsg{Width: 80, Height: 10}
)

func main() {
	cfg, err := loadConfig()
	check("load config", err)

	clientset, err := k8s.NewClientset(cfg.kubeconfig)
	check("create k8s clientset", err)

	prg := tea.NewProgram(models.Main(clientset))

	_, err = prg.Run()
	check("run program", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}

type config struct {
	kubeconfig string
}

func loadConfig() (config, error) {
	godotenv.Load()

	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return config{}, fmt.Errorf("get user home dir: %w", err)
		}

		kubeconfig = filepath.Join(homedir, ".kube", "config")
	}

	return config{kubeconfig: kubeconfig}, nil
}
