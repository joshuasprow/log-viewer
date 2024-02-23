package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuasprow/log-viewer/pkg"
)

type namespacesModel struct {
	model tea.Model
}

func newNamespacesModel() tea.Model {
	return namespacesModel{
		model: newListModel(pkg.Namespaces),
	}
}

func (n namespacesModel) Init() tea.Cmd {
	return nil
}

func (n namespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	n.model, cmd = n.model.Update(msg)
	return n, cmd
}

func (n namespacesModel) View() string {
	return n.model.View()
}
