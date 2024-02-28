package styles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var App = lipgloss.NewStyle().Width(80).Height(10)

var ListItem = struct {
	Normal   lipgloss.Style
	Selected lipgloss.Style
}{
	Normal:   lipgloss.NewStyle().PaddingLeft(4),
	Selected: lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")),
}

var List = struct {
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
