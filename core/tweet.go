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
