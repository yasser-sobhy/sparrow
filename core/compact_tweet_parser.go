package core

import (
	"bytes"
)

// CompactTweetParser parses tweets encoded using compact tweet format
// a compact tweet format is [code;arg1:arg1v:arg2:arg2v:argN:argNv;content]
// arguments and content are optional
// code, args, content must end with ;
type CompactTweetParser struct {
}

// Parse parses a tweet text and returns a Tweet type
func (parser *CompactTweetParser) Parse(message []byte, user *core.User) *core.Tweet {

	tweet := core.Tweet{}
	depth := 0 // depth of tweet part

	for i, char := range message {
		if char == ';' { //&& message[i-1] != '\\'
			if depth == 0 {
				tweet.Code = message[:i]
			} else if depth == 1 {
				tweet.Arguments = parser.processArguments(message[len(tweet.Code):i])
			} else if depth == 2 {
				tweet.Content = message[i:]
			}
			depth++
		}
	}

	tweet.Raw = message
	return &tweet
}

// Compile creates a raw tweet out of twwet parts
//example: char *args[]{tweet.code, tweet.argument, tweet.content}
func (parser *CompactTweetParser) Compile(code, arguments, content []byte) []byte {
	raw := [][]byte{code, arguments, content}
	return bytes.Join(raw, []byte{';'})
}

// Serialize serializes a tweet into a string
// different than raw, as the tweet may have been changed
// so, this function re-constructs the tweet's string
func (parser *CompactTweetParser) Serialize(tweet *core.Tweet) []byte {
	args := [][]byte{}

	for i, arg := range tweet.Arguments {
		args[i] = []byte{arg.Name..., ':', arg.Value}
	}

	return parser.Compile(tweet.Code, bytes.Join(args[:], []byte{':'}), tweet.Content)
}

func (parser *CompactTweetParser) processArguments(src []byte) []core.Arg {
	arguments := []core.Arg{}

	for i, char := range src {
		if char == ':' {
			arguments = append(arguments, core.Arg{Name: src[i:], Value: src[:i]})
		}
	}

	return arguments
}