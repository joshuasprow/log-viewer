package main

type viewKey string

const (
	namespacesKey viewKey = "namespaces"
	apisKey       viewKey = namespacesKey + ".apis"

	containersKey    viewKey = apisKey + ".containers"
	containerLogsKey viewKey = containersKey + ".logs"

	cronJobsKey          viewKey = apisKey + ".cronjobs"
	cronJobJobsKey       viewKey = cronJobsKey + ".jobs"
	cronJobContainersKey viewKey = cronJobJobsKey + ".containers"
	cronJobLogsKey       viewKey = cronJobContainersKey + ".logs"
)

func (k viewKey) FilterValue() string {
	return string(k)
}
