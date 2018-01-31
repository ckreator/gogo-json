package gogojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*type Iterator interface {
	next      func() string
	peek      func() string
	eof       func() bool
	suffocate func(string, ...string)
}*/

const basicString = `abc`
const newlineString = `ab
c`

func TestEmpty(t *testing.T) {
	assert := assert.New(t)
	iterator := MakeIterator("")

	assert.Equal(iterator.Peek(), "")
}

func TestBasic(t *testing.T) {
	assert := assert.New(t)
	iterator := MakeIterator(basicString)

	assert.Equal(iterator.Next(), "a")
	assert.Equal(iterator.GetRow(), uint64(2))
	assert.Equal(iterator.GetLine(), uint64(1))

	assert.Equal(iterator.Next(), "b")
	assert.Equal(iterator.GetRow(), uint64(3))
	assert.Equal(iterator.GetLine(), uint64(1))

	assert.Equal(iterator.Next(), "c")
	assert.Equal(iterator.GetRow(), uint64(4))
	assert.Equal(iterator.GetLine(), uint64(1))
}

func TestNewline(t *testing.T) {
	assert := assert.New(t)
	iterator := MakeIterator(newlineString)

	assert.Equal(iterator.Next(), "a")
	assert.Equal(iterator.GetRow(), uint64(2))
	assert.Equal(iterator.GetLine(), uint64(1))

	assert.Equal(iterator.Next(), "b")
	assert.Equal(iterator.GetRow(), uint64(3))
	assert.Equal(iterator.GetLine(), uint64(1))

	assert.Equal(iterator.Next(), "\n")
	assert.Equal(iterator.GetRow(), uint64(0))
	assert.Equal(iterator.GetLine(), uint64(2))

	assert.Equal(iterator.Next(), "c")
	assert.Equal(iterator.GetRow(), uint64(1))
	assert.Equal(iterator.GetLine(), uint64(2))

	assert.Equal(iterator.Next(), "")
}
