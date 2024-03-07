package models

type Titled interface {
	Title() string
}

type Described interface {
	Description() string
}
