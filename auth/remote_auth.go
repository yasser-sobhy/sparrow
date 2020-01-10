package auth

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yasser-sobhy/sparrow/core"
)

// BasicAuth allows users to login using a remote services
type RemoteAuth struct {
	sparrow     *core.Sparrow
	APIEndpoint string
	Token       string
	client      *resty.Client
}

// NewRemoteAuth creates a new RemoteAuth
func NewRemoteAuth(sparrow *core.Sparrow) *RemoteAuth {
	c := resty.New()
	c.R().SetAuthToken(viper.GetString("remote_auth.token"))
	return &RemoteAuth{
		sparrow:     sparrow,
		APIEndpoint: viper.GetString("remote_auth.api_endpoint"),
		Token:       viper.GetString("remote_auth.token"),
		client:      c,
	}
}

func (remoteAuth *RemoteAuth) Login(ws *core.WebSocket, token []byte) bool {

	resp, err := remoteAuth.client.R().SetBody(token).Post(remoteAuth.APIEndpoint)

	if err != nil {
		logrus.Error("RemoteAuth failed: status code {}, response {}", err)
		return false
	}

	if resp.StatusCode() != 200 {
		logrus.Error("RemoteAuth failed: status code {}, response {}", resp.StatusCode(), resp.Body())
		return false
	}

	logrus.Info("RemoteAuth, loggin in user {}", resp.Body())
	remoteAuth.sparrow.LogUserIn(token, core.USER, ws)
	return true
}
