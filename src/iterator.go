package gogojson

// Iterator is an interface that provides an easy way to interact with json input
// from whatever source
type Iterator interface {
	Next() string
	Peek() string
	Eof() bool
	// Suffocate(string, ...string)
}

type StringIterator struct {
	source  string
	current uint64
	line    uint64
	row     uint64
	length  uint64
}

func (iter *StringIterator) Next() string {
	if iter.Eof() {
		return ""
	}

	char := string(iter.source[iter.current])
	iter.current += 1

	if char == "\n" {
		iter.row = 0
		iter.line += 1
	} else {
		iter.row += 1
	}

	return char
}

func (iter *StringIterator) Peek() string {
	if iter.Eof() {
		return ""
	}

	return string(iter.source[iter.current])
}

func (iter *StringIterator) Eof() bool {
	return iter.length <= iter.current
}

func (iter *StringIterator) HasNext() bool {
	return !iter.Eof()
}

/*func (iter *StringIterator) Suffocate(msg string, additional ...string) {
	fmt.Printf("Error happened in iterator")
}*/

func (iter *StringIterator) GetLine() uint64 {
	return iter.line
}

func (iter *StringIterator) GetRow() uint64 {
	return iter.row
}

func MakeIterator(source string) *StringIterator {
	return &StringIterator{
		source:  source,
		current: 0,
		line:    1,
		row:     1,
		length:  uint64(len(source)),
	}
}
