package tui

import "github.com/charmbracelet/bubbles/list"

type Log string

func (l Log) FilterValue() string {
	return string(l)
}

func WrapLogs(logs []string) []list.Item {
	wrapped := make([]list.Item, len(logs))
	for i, l := range logs {
		wrapped[i] = Log(l)
	}
	return wrapped
}
