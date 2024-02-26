package models

type View string

const (
	NamespacesView View = "namespaces"
	ContainersView View = "containers"
	LogsView       View = "logs"
)
