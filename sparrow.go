package main

import "github.com/yasser-sobhy/sparrow/core"
import "github.com/yasser-sobhy/sparrow/middlewares"

func main() {
	var sparrow core.Sparrow
	sparrow.InstallMiddleware(middlewares.Debug{}, core.MiddlewareOptions{})
	sparrow.Run()
}
