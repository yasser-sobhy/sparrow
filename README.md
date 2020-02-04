**NOTE** Sparrow under heavy development and is broken(moving from C++ to Go)
```go
import ("github.com/yasser.sobhy/sparrow/core")

type Notification struct {
}

// send a notification to a user
func (n *Notification) Congregate(flock *core.Flock, options core.twitt) {
    process := func(ws *core.Conn, user core.User, tweet *core.Tweet) bool {
      // find target user
      if user := core.users.Get(tweet.Arguments[0].Value) {
        // deliver notification
        ws.send(tweet.Raw)
        return true
      }
      return false
    }

    flock.AddMiddleware(process, options)
}
```

## Tweets
Sparrow uses a very simple protocol to exchange messages (aka tweets). A Sparrow tweet consists of a code, a list of arguments, and a content separated by a semicolon.

Tweets arguments and content are optional, only code is required, for example:

```
n;arg1=4567;some text

n 			-> tweet code, specifies which target twitter callback to be run
arg1 		-> arguments
some text	-> tweet content
```

Sparrow.Tweet is a struct:

```go
  type Tweet struct {
    Code      []byte //target twitter callback
    Arguments []Arg  // key-value pairs argments
    Content   []byte // tweet's content
    Raw       []byte // message as received by uWebSockets
  }
```


## Users
Sparrow makes it very simple to add, remove and lookup users. Sparrow use a special [Trie](https://en.wikipedia.org/wiki/Trie) internally to store users:

```go
  // find user
  if user := sparrow.Users.Get("some_user_id") {
    user.write(...);
  }
```

Actually you should not need to add or remove users because Sparrow takes care of that, most of the time you will only need to find a user and do your logic. To allow users to login use one of:

`Sparrow.Twitters.BasicAuth` to allow users with ids stored in config file to login
`Sparrow.Twitters.RemoteAuth` to login users using external services
`Sparrow.Twitters.DevAuth` development-only twitter to allow any user to login without authentication


## Environment
Sparrow supports the four most common environments development, staging, and production.

```go

// set environment
  env.SetEnvironment("production");

// you can also check for environments and change your codes behavior based on that
  if(env.IsDevelopment()){
	  //do something...
  }

  // also availabe
  env.IsDevelopment()
  env.IsStaging()
  env.IsProduction()
```

## Remote events
You can use `Sparrow.RemoteEvent` to notify external services of events that has taken place in your code

```go
Sparrow.RemoteEvent.Trigger("login",  {user_id: ...});
// or
Sparrow.RemoteEvent.TriggerAsync("login",  {user_id: ...});
```
