package defaults

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var listItemStyles = struct {
	NormalTitle   lipgloss.Style
	SelectedTitle lipgloss.Style
	Description   lipgloss.Style
}{
	NormalTitle: lipgloss.NewStyle().PaddingLeft(4),
	SelectedTitle: lipgloss.
		NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("170")),
	Description: lipgloss.
		NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.Color("244")),
}

var ListStyles = struct {
	Help       lipgloss.Style
	NoItems    lipgloss.Style
	Pagination lipgloss.Style
	QuitText   lipgloss.Style
	Spinner    lipgloss.Style
	Title      lipgloss.Style
	TitleBar   lipgloss.Style
}{
	Help:       list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1),
	NoItems:    lipgloss.NewStyle().PaddingLeft(4),
	Pagination: lipgloss.NewStyle().PaddingLeft(4),
	QuitText:   lipgloss.NewStyle().Margin(1, 0, 2, 4),
	Spinner:    lipgloss.NewStyle(),
	Title:      lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
	TitleBar:   lipgloss.NewStyle().PaddingLeft(4),
}
