package middlewares

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yasser-sobhy/sparrow/core"
)

// This middleware log incomming tweets to a log file
type Log struct {
	Name         []byte
	Enabled      bool
	LogFullTweet bool
}

func NewLog() *Log {
	viper.SetDefault("enabled", true)
	viper.SetDefault("log_full_tweet", false)

	return &Log{
		Name:         []byte("Log"),
		Enabled:      viper.GetBool("enabled"),
		LogFullTweet: viper.GetBool("log_full_tweet"),
	}
}

func (log *Log) Congregate(flock *core.Flock, options core.MiddlewareOptions) {
	// TODO use a dedicated logger here not Sparrow::Logger
	process := func(ws *core.Conn, user core.User, tweet *core.Tweet) bool {
		if log.LogFullTweet {
			logrus.Info("New Tweet: {}{}", user.ID, tweet.Raw)
		} else {
			logrus.Info("New Tweet: {}{}", user.ID, tweet.Code)
		}
		return true
	}

	if log.Enabled {
		flock.AddMiddleware(process, options)
	}
}
