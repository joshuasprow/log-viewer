package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/joshuasprow/log-viewer/k8s"
)

type Container struct {
	k8s.Container
}

func (c Container) FilterValue() string {
	return fmt.Sprintf("%s.%s.%s", c.Namespace, c.Pod, c.Name)
}

func WrapContainers(containers []k8s.Container) []list.Item {
	wrapped := make([]list.Item, len(containers))
	for i, c := range containers {
		wrapped[i] = Container{c}
	}
	return wrapped
}
