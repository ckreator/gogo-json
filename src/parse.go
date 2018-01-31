package gogojson

import (
	"fmt"
)

func parseRecursive(tokens []*Token, host map[string]interface{}) ([]*Token, map[string]interface{}) {
	var key string
	var next *Token

	for len(tokens) > 0 {
		key, tokens = assertString(tokens)
		tokens = skipPunctuation(tokens, ":")
		next, tokens = shift(tokens)

		if next.Type == STRING || next.Type == NUM || next.Type == BOOL || next.Type == NULL {
			host[key] = next.Value
		} else if next.Value == "{" {
			tokens, host[key] = parseRecursive(tokens, newMap())
		}

		if len(tokens) > 0 && tokens[0].Type == PUNC && tokens[0].Value == "}" {
			tokens = skipPunctuation(tokens, "}")
			break
		}
		tokens = skipPunctuation(tokens, ",")
	}

	return tokens, host
}

func shift(tokens []*Token) (*Token, []*Token) {
	return tokens[0], tokens[1:]
}

func assertString(tokens []*Token) (string, []*Token) {
	key, tokens := shift(tokens)
	if key.Type == STRING {
		return key.Value.(string), tokens
	}

	panic(fmt.Sprintf("Expected string key, instead got a %s", key.Type))
}

func newMap() map[string]interface{} {
	return make(map[string]interface{})
}

// Parse is a function that takes in a list of Token Pointers and returns a
// generic map type for the json object
func Parse(input []*Token) map[string]interface{} {
	// skip curly bracket
	tokens := skipPunctuation(input, "{")
	_, out := parseRecursive(tokens, newMap())

	return out
}

func skipPunctuation(tokens []*Token, when string) []*Token {
	t, tokens := tokens[0], tokens[1:]
	if t.Type == PUNC && t.Value == when {
		return tokens
	}

	panic(fmt.Sprintf("Expected punctuation with value '%s', instead got: '%s'", when, t.Value.(string)))
}
