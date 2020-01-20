package core

type Arg struct {
	Name  []byte
	Value []byte
}

type Tweet struct {
	Code      []byte //target twitter callback
	Arguments []Arg  // key-value pairs argments
	Content   []byte // tweet's content
	Raw       []byte // message as received by ws
}

func (tweet Tweet) Valid() bool {
	// code is the only required attribute. arguments and content are optional
	return tweet.Code != nil
}

// duplicte tweet data
// because tweets are deleted after running a twitter
// these functions may be used to duplicate tweet data to keep it for future use
func (tweet Tweet) Dup() *Tweet {
	t := Tweet{}

	t.Code = make([]byte, len(tweet.Code))
	t.Arguments = make([]Arg, len(tweet.Arguments))
	t.Content = make([]byte, len(tweet.Content))
	t.Raw = make([]byte, len(tweet.Raw))

	copy(t.Code, tweet.Code)
	copy(t.Arguments, tweet.Arguments)
	copy(t.Content, tweet.Content)
	copy(t.Raw, tweet.Raw)

	return &t
}
