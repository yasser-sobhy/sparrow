package core

import (
	"bytes"
)

// CompactTweetParser parses tweets encoded using compact tweet format
// a compact tweet format is [code;arg1=arg1v,arg2=arg2v,argN=argNv;content]
// arguments and content are optional
// code, args, content must end with ;
type CompactTweetParser struct {
}

// Parse parses a tweet text and returns a Tweet type
func (parser *CompactTweetParser) Parse(message []byte) *Tweet {

	tweet := Tweet{}
	depth := 0   // depth of tweet part
	argsEnd := 0 // index of args ;

	for i, char := range message {
		if char == ';' { //&& message[i-1] != '\\'
			if depth == 0 {
				tweet.Code = string(message[:i])
			} else if depth == 1 {
				argsEnd = i + 1
				tweet.Args = parser.processArguments(message[len(tweet.Code)+1 : i])
			} else if depth == 2 {
				tweet.Content = message[argsEnd:i]
			}
			depth++
		}
	}

	tweet.Raw = message
	return &tweet
}

// Compile creates a raw tweet out of twwet parts
//example: char *args[]{tweet.code, tweet.argument, tweet.content}
func (parser *CompactTweetParser) Compile(code, arguments string, content []byte) []byte {
	raw := [][]byte{[]byte(code), []byte(arguments), content}
	return bytes.Join(raw, []byte{';'})
}

// Serialize serializes a tweet into a string
// different than raw, as the tweet may have been changed
// so, this function re-constructs the tweet's string
func (parser *CompactTweetParser) Serialize(tweet *Tweet) []byte {
	args := ""

	for k, v := range tweet.Args {
		args += k + ":" + v
	}

	return parser.Compile(tweet.Code, args, tweet.Content)
}

func (parser *CompactTweetParser) processArguments(src []byte) map[string]string {
	arguments := map[string]string{}

	sep := 0
	beg := 0
	for i, char := range src {

		if char == '=' {
			sep = i
			continue
		}

		if char == ',' {
			arguments[string(src[beg:sep])] = string(src[sep+1 : i])
			beg = i + 1
		}
	}

	// read last arg if there was a separator after the last ,
	if sep > beg {
		arguments[string(src[beg:sep])] = string(src[sep+1:])
	}

	return arguments
}
