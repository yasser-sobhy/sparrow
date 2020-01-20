package middlewares

import "github.com/spf13/viper"
import "github.com/yasser-sobhy/sparrow/core"

// Cap closes websockets that send a tweet with length more than the configured length
// more caps may be added alter
// Configs can be changes using toml file under the key sparrow.tiwtters.cap
//      [sparrow.tiwtters.cap]
//      MaxMessageSize = 1024
// or using:
//      Sparrow.Twitters.Cap.MaxMessageSize = 1024
type Cap struct {
	// public configurations
	Name           []byte
	MaxMessageSize int
}

// NewCap Create a new Cap middleware
func NewCap() *Cap {
	viper.SetDefault("max_message_size", 1024)

	return &Cap{
		Name:           []byte("Cap"),
		MaxMessageSize: viper.GetInt("max_message_size"),
	}
}

// Congregate registers cap middleware
func (cap *Cap) Congregate(flock *core.Flock, options core.MiddlewareOptions) {
	process := func(ws *core.Conn, user *core.User, tweet *core.Tweet) bool {
		if len(tweet.Raw) > cap.MaxMessageSize {
			ws.Close()   // may not close, just return false to stop processing messages
			return false //interrput processing
		}
		return true
	}
	flock.AddMiddleware(process, options)
}
