package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var defaultSize = tea.WindowSizeMsg{Width: 80, Height: 10}

var listItemStyles = struct {
	Normal   lipgloss.Style
	Selected lipgloss.Style
}{
	Normal:   lipgloss.NewStyle().PaddingLeft(4),
	Selected: lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")),
}

var listStyles = struct {
	Pagination lipgloss.Style
	Help       lipgloss.Style
	Title      lipgloss.Style
	TitleBar   lipgloss.Style
	QuitText   lipgloss.Style
}{
	Pagination: lipgloss.NewStyle().PaddingLeft(4),
	Help:       list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1),
	Title:      lipgloss.NewStyle().MarginLeft(4).Foreground(lipgloss.Color("205")),
	TitleBar:   lipgloss.NewStyle(),
	QuitText:   lipgloss.NewStyle().Margin(1, 0, 2, 4),
}
