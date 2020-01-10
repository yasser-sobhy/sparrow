package events

import "github.com/spf13/viper"
import "github.com/yasser-sobhy/sparrow/net"
import "github.com/yasser-sobhy/sparrow/core"

// Login is used to notify external api's of user logins
type Login struct {
	Name        string
	APIEndpoint string
}

func NewLogin() Login {
	return Login{
		Name:        "Login",
		APIEndpoint: viper.GetString("api_endpoint"),
	}
}

// Congregate registers middleware functions to flock
func (login *Login) Congregate(flock *core.Flock, options core.MiddlewareOptions) {
	process := func(ws *core.WebSocket, user *core.User, tweet *core.Tweet) bool {
		// trigger a login event
		net.NewRemoteEvent(login.APIEndpoint).Trigger([]byte("login"), tweet.Arguments[0].Value)
		return true
	}

	flock.AddMiddleware(process, options)
}
