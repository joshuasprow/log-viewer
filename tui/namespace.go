package tui

import "github.com/charmbracelet/bubbles/list"

type Namespace string

func (n Namespace) FilterValue() string {
	return string(n)
}

func WrapNamespaces(namespaces []string) []list.Item {
	wrapped := make([]list.Item, len(namespaces))
	for i, n := range namespaces {
		wrapped[i] = Namespace(n)
	}
	return wrapped
}
