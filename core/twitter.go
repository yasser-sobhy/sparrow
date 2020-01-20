package core

type TwitterOptions struct {
	Code  byte
	Scope Scope
	Post  bool
	Async bool
}

type Twitter interface {
	// this twitter's name
	Name() string
	// a function to register the twitter's handler functions
	Congregate(flock *Flock)
}
