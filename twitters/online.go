package twitters

import ("github.com/yasser.sobhy/sparrow/core")

type Online struct {
}

func (twitter *Twitter) Congregate(flock *core.Flock, options core.twitt) {
    process := func(ws *core.WebSocket, user *core.User, tweet *core.Tweet) bool {
        if user := users.Get(tweet.Arguments[0].Value) {
            ws.send(tweet.Raw)
            return true
        }
        return false
    }

    flock.AddMiddleware(process, options)
}