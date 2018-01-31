package gogojson

import (
	"fmt"
	"regexp"
	"strconv"
)

type Token struct {
	Type  string
	Value interface{}
}

const PUNC = "PUNC"
const NUM = "NUM"
const STRING = "STRING"
const BOOL = "BOOL"
const NULL = "NULL"

// define all regexes for parsing
var punctuationRegex = regexp.MustCompile("[{}:,\\[\\]]")

// check for true, false and null
var identifierRegex = regexp.MustCompile("[truefalsn]")
var numberRegex = regexp.MustCompile("[0-9]")
var stringInitRegex = regexp.MustCompile("\"")
var stringBodyRegex = regexp.MustCompile("[^\"]")
var whitespaceRegex = regexp.MustCompile("[\n\t\r ]")

func Tokenize(iter *StringIterator) []*Token {
	tokens := make([]*Token, 0)
	nextToken := func(next string) *Token {
		if isPunctuation(next) {
			return iter.ParsePunctuation(next)
		} else if isStringInit(next) {
			return iter.ParseString(next)
		} else if isIdentifier(next) {
			return iter.ParseIdentifier(next)
		} else if isNumber(next) {
			return iter.ParseNumber(next)
		}

		panic(fmt.Sprintf("Unexpected character type: '%s'", next))
	}

	for iter.HasNext() {
		// skip whitespace as we don't care about it
		iter.SkipWhitespace()

		next := iter.Next()
		tokens = append(tokens, nextToken(next))
	}

	return tokens
}

// Character classification
func isPunctuation(check string) bool {
	return punctuationRegex.MatchString(check)
}

func isStringInit(check string) bool {
	return stringInitRegex.MatchString(check)
}

func isStringBody(check string) bool {
	return !stringInitRegex.MatchString(check)
}

func isIdentifier(check string) bool {
	return identifierRegex.MatchString(check)
}

func isNumber(check string) bool {
	return numberRegex.MatchString(check)
}

func isWhitespace(check string) bool {
	return whitespaceRegex.MatchString(check)
}

// iterator parse methods
func (json *StringIterator) ParsePunctuation(init string) *Token {
	return &Token{
		Value: init,
		Type:  PUNC,
	}
}

func (json *StringIterator) ParseString(init string) *Token {
	str := ""
	for isStringBody(json.Peek()) {
		str += json.Next()
	}

	// skip closing
	json.Next()

	return &Token{
		Value: str,
		Type:  STRING,
	}
}

func (json *StringIterator) ParseNumber(init string) *Token {
	for isNumber(json.Peek()) {
		init += json.Next()
	}

	num, _ := strconv.ParseFloat(init, 32)

	// TODO: make it support floats
	return &Token{
		Value: num,
		Type:  NUM,
	}
}

func (json *StringIterator) ParseIdentifier(init string) *Token {
	for isIdentifier(json.Peek()) {
		init += json.Next()
	}

	if init == "true" || init == "false" {
		return &Token{
			Value: init == "true",
			Type:  BOOL,
		}
	}

	return &Token{
		Value: nil,
		Type:  NULL,
	}
}

// TODO: improve this maybe?
func (json *StringIterator) SkipWhitespace() {
	for isWhitespace(json.Peek()) {
		json.Next()
	}
}
