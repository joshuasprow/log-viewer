package main

import "github.com/charmbracelet/lipgloss"

func renderTitle(path ...string) string {
	var title string
	for i, p := range path {
		if i > 0 {
			title += " > "
		}
		color := "#FF00FF"
		if i == len(path)-1 {
			color = "#00FF00"
		}
		title += lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Render(p)
	}
	return title
}
