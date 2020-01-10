package net

import (
	"github.com/go-resty/resty/v2"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RemoteEvent struct {
	APIEndpoint string
	Token       string

	client   *resty.Client
	enqueuer *work.Enqueuer
}

func NewRemoteEvent(apiEndpoint string) RemoteEvent {
	c := resty.New()
	c.R().SetAuthToken(viper.GetString("remote_events.token"))

	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379")
		},
	}
	// Make an enqueuer with a particular namespace
	var enqueuer = work.NewEnqueuer("my_app_namespace", redisPool)

	if apiEndpoint != "" {
		apiEndpoint = viper.GetString("remote_events.api_endpoint")
	}

	return RemoteEvent{
		APIEndpoint: apiEndpoint,
		Token:       viper.GetString("remote_events.token"),
		client:      c,
		enqueuer:    enqueuer,
	}
}

func (remoteEvent RemoteEvent) Trigger(event []byte, payload []byte) {
	_, err := remoteEvent.enqueuer.Enqueue("remote_event", work.Q{"event": event, "payload": payload})
	if err != nil {
		logrus.Error("RemoteEvent failed to enqueue job {}", err)
	}
}

func (remoteEvent RemoteEvent) TriggerAsync(event string, payload []byte) {
	go func() {
		resp, err := remoteEvent.client.R().SetBody(payload).Post(remoteEvent.APIEndpoint + "/" + event)

		if err != nil {
			logrus.Error("RemoteEvent failed: status code {}, response {}", err)
		}

		if resp.StatusCode() != 200 {
			logrus.Error("RemoteEvent failed: status code {}, response {}", resp.StatusCode(), resp.Body())
		}

		logrus.Info("RemoteEvent, loggin in user {}", resp.Body())
	}()
}
