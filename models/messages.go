package models

type ErrMsg struct{ Err error }

func (e ErrMsg) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

type NamespacesMsg []string

type ContainersMsg []containerListItem
