package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yasser-sobhy/sparrow/core"
	"github.com/yasser-sobhy/sparrow/middlewares"
)

func main() {
	var sparrow *core.Sparrow = core.NewSparrow()
	sparrow.InstallMiddleware(&middlewares.Debug{}, core.MiddlewareOptions{})

	sparrow.Use(func(c *core.Conn, u core.User, t *core.Tweet) bool {
		logrus.Info(string(t.Code), string(t.Content))
		c.WriteMessage([]byte("Hello Sparrow!"))
		return true
	}, core.TwitterOptions{Code: "0", Scope: core.NONE})
	sparrow.Run()
}
