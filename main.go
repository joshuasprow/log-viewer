package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func main() {
	d, err := readModelData()
	check("read model data", err)

	p := tea.NewProgram(newModel(d))

	_, err = p.Run()
	check("run program", err)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}

type flags struct {
	kubeconfig string
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

type podData struct {
	Name string   `json:"name"`
	Logs []string `json:"logs"`
}

type namespaceData struct {
	Name string    `json:"name"`
	Pods []podData `json:"pods"`
}

type modelData struct {
	namespaces []namespaceData
}

func readModelData() (modelData, error) {
	d, err := os.ReadFile("tmp/data.json")
	if err != nil {
		return modelData{}, fmt.Errorf("read file: %w", err)
	}

	namespaces := []namespaceData{}

	err = json.Unmarshal(d, &namespaces)
	if err != nil {
		return modelData{}, fmt.Errorf("unmarshal data: %w", err)
	}

	return modelData{namespaces: namespaces}, nil
}

type model struct {
	data modelData
}

func newModel(data modelData) model {
	return model{data: data}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return model{}, tea.Quit
		}
	}

	return model{}, nil
}

func (m model) View() string {
	return "hello"
}
