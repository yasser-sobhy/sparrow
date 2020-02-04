package auth

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yasser-sobhy/sparrow/core"
)

// BasicAuth allows users to login using a username and password
// RemoteAuth
type BasicAuth struct {
	sparrow *core.Sparrow
	// users allowed to login
	// read from configuration file or environment variables
	userIDs map[string]string
}

// NewBasicAuth creates a new BasicAuth. Reads users data from config files
func NewBasicAuth(sparrow *core.Sparrow) BasicAuth {
	userIDs := viper.GetStringMapString("basic_auth.users")
	return BasicAuth{userIDs: userIDs, sparrow: sparrow}
}

// Login processes the incoming tweet and allows user to login if user name and password
// found in UserIDs map
func (basicAuth *BasicAuth) Login(ws *core.Conn, user core.User, tweet *core.Tweet) {
	userInfo := tweet.Arguments[0]
	if info, ok := basicAuth.userIDs[string(userInfo.Name)]; ok && info == string(userInfo.Value) {
		logrus.Info("BasicAuth, loggin in user {}", userInfo.Name)
		basicAuth.sparrow.LogUserIn(tweet.Arguments[0].Value, core.USER, ws)
	} else {
		logrus.Warn("BasicAuth, user not fount {}", userInfo.Name)
		//ws.close()
	}
}
