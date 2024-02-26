package models

type ViewKey string

const (
	NamespacesView ViewKey = "namespaces"
	ContainersView ViewKey = "containers"
	LogsView       ViewKey = "logs"
)
