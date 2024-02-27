package models

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var appStyles = lipgloss.NewStyle().Width(80).Height(10).Padding(0, 2)

var listItemStyles = struct {
	Normal   lipgloss.Style
	Selected lipgloss.Style
}{
	Normal:   lipgloss.NewStyle().PaddingLeft(4),
	Selected: lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")),
}

var listStyles = struct {
	Filter     lipgloss.Style
	Help       lipgloss.Style
	Pagination lipgloss.Style
	QuitText   lipgloss.Style
	Title      lipgloss.Style
	TitleBar   lipgloss.Style
}{
	Help:       list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1),
	Pagination: lipgloss.NewStyle().PaddingLeft(4),
	QuitText:   lipgloss.NewStyle().Margin(1, 0, 2, 4),
	Title:      lipgloss.NewStyle().MarginLeft(4).Foreground(lipgloss.Color("205")),
	TitleBar:   lipgloss.NewStyle(),
}
