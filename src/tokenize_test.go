package gogojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var json = `{
  "hello": 123,
  "how": {},
  "lies": true
}`

func TestTokenize(t *testing.T) {
	iterator := MakeIterator(json)
	tokens := Tokenize(iterator)

	assert := assert.New(t)

	assert.Equal(PUNC, tokens[0].Type, "Type is not PUNC")
	assert.Equal(STRING, tokens[1].Type, "Type is not STRING")
	assert.Equal("hello", tokens[1].Value, "Value is not \"hello\"")
	assert.Equal(float64(123), tokens[3].Value, "Value is not 123")

	/*for _, tok := range tokens {
		fmt.Println(tok.Type, tok.Value)
	}*/
}

func TestTokenizeNil(t *testing.T) {
	assert := assert.New(t)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(r, "Unexpected character type: '!'")
		}
	}()

	iterator := MakeIterator(`!`)
	Tokenize(iterator)
}
