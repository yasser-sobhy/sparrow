package middlewares

import (
	"github.com/sirupsen/logrus"
	"github.com/yasser-sobhy/sparrow/core"
	"github.com/yasser-sobhy/sparrow/env"
)

// This middleware prints the tweet to STDOUT
// it will can also provide more functionality for development environment ONLY
type Debug struct {
	Name []byte
}

func NewDebug() *Debug {
	if env.IsDevelopment() {
		logrus.Warn("Debug: Warning, Debug is being used in wrong environment")
	}
	return &Debug{[]byte("Debug")}
}

func (debug *Debug) Congregate(flock *core.Flock, options core.MiddlewareOptions) {
	process := func(ws *core.Conn, user *core.User, tweet *core.Tweet) bool {
		logrus.Info("New Tweet: {}", tweet.Raw)
		return true
	}

	if env.IsDevelopment() {
		flock.AddMiddleware(process, options)
	}
}
