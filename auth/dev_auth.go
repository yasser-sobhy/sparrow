package auth

import (
	"github.com/sirupsen/logrus"
	"github.com/yasser-sobhy/sparrow/core"
	"github.com/yasser-sobhy/sparrow/env"
)

// DevAuth is a development only twitter
// this twitter allows any users to login without authentecation
type DevAuth struct {
	sparrow *core.Sparrow
}

// NewDevAuth creates a new DevAuth
func NewDevAuth(sparrow *core.Sparrow) DevAuth {
	if env.IsDevelopment() {
		return DevAuth{sparrow: sparrow}
	}

	panic("DevAuth: Warning, DevAuth is being used in wrong environment")
}

// Login processes the incoming tweet and allows user to login
func (devAuth *DevAuth) Login(ws *core.WebSocket, token []byte) bool {
	logrus.Info("DevAuth, loggin in user {}", token)
	devAuth.sparrow.LogUserIn(token, core.USER, ws)
	return true
}
