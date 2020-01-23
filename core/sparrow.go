package core

import (
	"github.com/dghubble/trie"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/evio"
)

// Sparrow server
type Sparrow struct {
	// contains twitters, followers, leaders, middlewares, and postwares
	Flock       Flock
	TweetParser TweetParser
	Users       *trie.PathTrie
	Channels    *trie.PathTrie
}

// New creates a new Sparrow instance
func NewSparrow() *Sparrow {
	return &Sparrow{
		TweetParser: &CompactTweetParser{},
		Flock:       NewFlock(),
		Users:       trie.NewPathTrie(),
		Channels:    trie.NewPathTrie(),
	}
}

func (sparrow *Sparrow) Install(twitter Twitter) {
	twitter.Congregate(&sparrow.Flock)
}
func (sparrow *Sparrow) Use(twitter TweetHandler, options TwitterOptions) {
	sparrow.Flock.Add(twitter, options)
}

// middleware run order:
// 1- middlewares
// 2- twitter-specific middlewares
// 3- twitter
// 4- post middlewares
// 5- twitter-specific post middlewares
func (sparrow *Sparrow) InstallMiddleware(middleware Middleware, options MiddlewareOptions) {
	middleware.Congregate(&sparrow.Flock, options)
}

// callback middlewares (lambda instead of classes)
func (sparrow *Sparrow) UseMiddleware(middleware TweetHandler, options MiddlewareOptions) {
	sparrow.Flock.AddMiddleware(middleware, options)
}

// install a single middleware to be run when a new ws is connected
func (sparrow *Sparrow) OnConnection(middleware OnConnectionHandler) {
	sparrow.Flock.AddOnConnection(middleware)
}

// install a single middleware to be run when a ws is disconnected
func (sparrow *Sparrow) OnDisconnection(middleware OnDisconnectionHandler) {
	sparrow.Flock.AddOnDisconnection(middleware)
}

func (sparrow *Sparrow) LogUserIn(id []byte, scope Scope, ws *Conn) bool {
	if sparrow.Users.Put(string(id), ws) {
		//ws.setUserData(new User{id, scope})
		return true
	}
	return false
}

func (sparrow *Sparrow) LogUserOut(id []byte, ws Conn) bool {
	success := sparrow.Users.Delete(string(id))
	//ws.setUserData(nullptr)
	return success
}

//Run starts listening for incoming messages
func (sparrow *Sparrow) Run() {

	nano := NanoWebsocket{}
	nano.OnOpen = func(c *Conn, handshake *ws.Handshake) {
		for _, middleware := range sparrow.Flock.OnConnectionMiddlewares {
			middleware(c, handshake)
		}
		return
	}

	nano.OnClose = func(c *Conn, err error) (action evio.Action) {
		user, _ := c.Context().(User)

		for _, middleware := range sparrow.Flock.OnDisconnectionMiddlewares {
			middleware(c, &user)
		}
		return
	}

	nano.OnMessage = func(c *Conn, message wsutil.Message) (out []byte, action evio.Action) {
		user, userOk := c.Context().(*User)
		tweet := sparrow.TweetParser.Parse(message.Payload)
		scope := NONE

		if userOk {
			scope = user.Scope
		}

		if tweet.Valid() {
			middlewares := sparrow.Flock.GetMiddlewares(scope)
			twitter, twitterOk := sparrow.Flock.Get(tweet.Code, scope)
			postMiddlewares := sparrow.Flock.GePosttMiddlewares(scope)

			if twitterOk {
				proceed := true

				// run middlware, any middlware that return false will interrupt execution, tiwtter will not run
				for _, middleware := range middlewares {
					proceed = proceed && middleware(c, user, tweet)
				}

				if proceed {
					twitter(c, user, tweet)
				}

				// run postware
				for _, middleware := range postMiddlewares {
					middleware(c, user, tweet)
				}
			}
			//else log('twitter not found', message)
		} else {
			// log('recieved invalid tweet', message)
		}

		//delete tweet
		return
	}

	viper.SetDefault("sparrow.port", 9001)
	port := viper.GetString("sparrow.port")
	if err := nano.Serve("tcp://localhost:" + port); err != nil {
		panic(err.Error())
	}

	logrus.Info("Sparrow, listening on port {}", port)
}
