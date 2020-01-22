package core

import "github.com/gobwas/ws"

// Conn, user,   tweet,    twitter,      success: the result of the ran twitter, for post middlewares
type TweetHandler func(*Conn, *User, *Tweet) bool
type OnConnectionHandler func(*Conn, *ws.Handshake) bool
type OnDisconnectionHandler func(*Conn, *User) bool

type Flock struct {
	twitters map[Scope]map[string]TweetHandler

	middlewares     map[Scope][]TweetHandler
	postMiddlewares map[Scope][]TweetHandler

	OnConnectionMiddlewares    []OnConnectionHandler
	OnDisconnectionMiddlewares []OnDisconnectionHandler
}

func NewFlock() Flock {
	return Flock{
		twitters: map[Scope]map[string]TweetHandler{
			ANY:   map[string]TweetHandler{},
			NONE:  map[string]TweetHandler{},
			USER:  map[string]TweetHandler{},
			ADMIN: map[string]TweetHandler{},
		},

		middlewares:     map[Scope][]TweetHandler{},
		postMiddlewares: map[Scope][]TweetHandler{},

		OnConnectionMiddlewares:    []OnConnectionHandler{},
		OnDisconnectionMiddlewares: []OnDisconnectionHandler{},
	}
}

// tiwtter callbacks. These are the actual callbacks that will process tweets
// they're not middlewares
func (flock *Flock) Add(twitter TweetHandler, options TwitterOptions) {
	// if twitter found
	flock.twitters[options.Scope][options.Code] = twitter
}

func (flock *Flock) AddMany(twitters []TweetHandler, options TwitterOptions) {
	for _, twitter := range twitters {
		flock.Add(twitter, options)
	}
}

// retrieve a twitter
func (flock *Flock) Get(code string, scope Scope) (TweetHandler, bool) {
	// TDOD: handle duplicate twitter
	if v, ok := flock.twitters[scope]; ok {
		return v[code], true
	}
	return nil, false
}

// middlewares
func (flock *Flock) AddMiddleware(twitter TweetHandler, options MiddlewareOptions) {
	if options.Post {
		flock.postMiddlewares[options.Scope] = append(flock.postMiddlewares[options.Scope], twitter)
	} else {
		flock.middlewares[options.Scope] = append(flock.middlewares[options.Scope], twitter)
	}
}

// middlewares with Scope.Any should be returned here, even if scope was diffferent
func (flock *Flock) GetMiddlewares(scope Scope) []TweetHandler {
	return flock.postMiddlewares[scope]
}

// middlewares with Scope.Any should be returned here, even if scope was diffferent
func (flock *Flock) GePosttMiddlewares(scope Scope) []TweetHandler {
	return flock.middlewares[scope]
}

// install a single middleware to be run when a new ws is connected
func (flock *Flock) AddOnConnection(middleware OnConnectionHandler) {
	flock.OnConnectionMiddlewares = append(flock.OnConnectionMiddlewares, middleware)
}

// install a single middleware to be run when a ws is disconnected
func (flock *Flock) AddOnDisconnection(middleware OnDisconnectionHandler) {
	flock.OnDisconnectionMiddlewares = append(flock.OnDisconnectionMiddlewares, middleware)
}
