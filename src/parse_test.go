package gogojson

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var parseJSON = `{
  "name": "Peter",
  "age": 42,
  "male": true,
  "config": {
    "surname": "Griffin",
    "secrets": null
  }
}`

func TestParse(t *testing.T) {
	iterator := MakeIterator(parseJSON)
	tokens := Tokenize(iterator)
	parsed := Parse(tokens)

	assert := assert.New(t)
	name := parsed["name"]
	age := parsed["age"]
	assert.Equal("Peter", name, "Age is not 42")
	assert.Equal(float64(42), age, "Age is not 42")

	fmt.Println("Parsed:", parsed)
}

func TestParsePanicKey(t *testing.T) {
	assert := assert.New(t)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(r, "Expected string key, instead got a NUM")
		}
	}()

	iterator := MakeIterator(`{123: 123}`)
	tokens := Tokenize(iterator)
	Parse(tokens)

}

func TestParsePanicPunc(t *testing.T) {
	assert := assert.New(t)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(r, "Expected punctuation with value ':', instead got: ','")
		}
	}()

	iterator := MakeIterator(`{"hello", 123}`)
	tokens := Tokenize(iterator)
	Parse(tokens)

}
