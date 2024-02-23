package cli

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/pkg"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type namespacesModel struct {
	clientset *kubernetes.Clientset
	model     tea.Model
}

func newNamespacesModel(clientset *kubernetes.Clientset) tea.Model {
	return namespacesModel{
		clientset: clientset,
		model:     newListModel(pkg.Namespaces),
	}
}

func (n namespacesModel) Init() tea.Cmd {
	return func() tea.Msg {
		list, err := n.clientset.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
		if err != nil {
			return fmt.Errorf("list namespaces: %w", err)
		}

		namespaces := []string{}

		for _, item := range list.Items {
			namespaces = append(namespaces, item.Name)
		}

		return namespaces
	}
}

func (n namespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	n.model, cmd = n.model.Update(msg)

	switch msg := msg.(type) {
	case error:
		panic(fmt.Errorf("namespacesModel.Update: %w", msg))
	case []string:
		n.model = newListModel(msg)
	}
	return n, cmd
}

func (n namespacesModel) View() string {
	return n.model.View()
}
