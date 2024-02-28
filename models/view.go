package models

type ViewKey string

const (
	NamespacesView ViewKey = "namespaces"
	ContainersView ViewKey = "containers"
	CronJobsView   ViewKey = "cronJobs"
	LogsView       ViewKey = "logs"
)
