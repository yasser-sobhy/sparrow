package core

type User interface {
	ID() string
	Scope() Scope
}
