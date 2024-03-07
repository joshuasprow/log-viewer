package tui

import "github.com/charmbracelet/bubbles/list"

type Api string

func (a Api) FilterValue() string {
	return string(a)
}

const (
	ContainersApi Api = "containers"
	CronJobsApi   Api = "cron jobs"
)

func GetApis() []list.Item {
	return []list.Item{
		ContainersApi,
		CronJobsApi,
	}
}
