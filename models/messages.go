package models

import "github.com/joshuasprow/log-viewer/k8s"

type ErrMsg struct{ Err error }

func (e ErrMsg) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

type NamespacesMsg []string

type ContainersMsg []k8s.Container

type LogsMsg []string
