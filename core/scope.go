package core

type Scope int32

const (
	NONE  Scope = 0 // before loging in
	ANY   Scope = 1 // all scopes
	USER  Scope = 2 // a user
	ADMIN Scope = 3 // an admin
)
