package core

type Scope int32

const (
	// rename ANY to ALL, or NONE to ANONYMOUS ?
	ANY   Scope = 0 // all scopes
	NONE  Scope = 1 // before loging in
	USER  Scope = 2 // a user
	ADMIN Scope = 3 // an admin
)
