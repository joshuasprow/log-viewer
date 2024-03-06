package tui

import "github.com/charmbracelet/bubbles/list"

type Api string

func (a Api) FilterValue() string {
	return string(a)
}

const (
	ContainersApi Api = "Containers"
	CronJobsApi   Api = "CronJobs"
)

func GetApis() []list.Item {
	return []list.Item{
		ContainersApi,
		CronJobsApi,
	}
}
