package core

// MiddlewareOptions can be used to change a middleware's scope or run it before/after twitter
type MiddlewareOptions struct {
	Scope Scope
	Post  bool
	Async bool
}

// Middleware sparrow middlewares should implement this interface
type Middleware interface {
	// register middleware
	Congregate(flock *Flock, options MiddlewareOptions)
}
