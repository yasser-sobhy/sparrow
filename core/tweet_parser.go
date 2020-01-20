package core

// TweetParser an interface for twwet parsers
type TweetParser interface {
	Parse(message []byte) *Tweet
}
