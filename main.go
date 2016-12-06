package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// just a fun little iterator
/*type Iterator interface {
	next      func() string
	peek      func() string
	eof       func() bool
	suffocate func(string, ...string)
}*/

// there are just a few types of tokens
// - punctuation :, { or }
// - string
// - number
// - object
// - array
// - boolean (true or false)
// - null

type Token struct {
	tok_type string
	value    string
}

type InputStream struct {
	next      func() string
	peek      func() string
	eof       func() bool
	suffocate func(string, ...string)
}

type Tokenizer struct {
	next      func() *Token
	peek      func() *Token
	eof       func() bool
	suffocate func(string, ...string)
}

func make_iterator(input string) *InputStream {
	row := 0
	col := 0
	pos := 0
	ret := new(InputStream)
	splitted := strings.Split(input, "")

	next := func() string {
		c := splitted[pos]
		pos++
		if c == "\n" {
			col = 0
			row++
		} else {
			col++
		}
		return c
	}

	eof := func() bool {
		return pos >= len(splitted)
	}

	peek := func() string {
		if eof() {
			return ""
		} else {
			return splitted[pos]
		}
	}

	suffocate := func(msg string, list ...string) {
		log.Fatal(msg)
	}

	ret.next = next
	ret.peek = peek
	ret.eof = eof
	ret.suffocate = suffocate
	return ret
}

func tokenizer(input *InputStream) *Tokenizer {
	ret := new(Tokenizer)
	//i := 0
	var curr_tok *Token

	// type parsers:
	parse_string := func(bracket string) *Token {
		out := ""
		var val string
		for !input.eof() {
			val = input.next()
			//fmt.Println("VAL STR: ", val, " | ", out)
			// handle escapes first
			if val == "\\" {
				out += val + input.next()
			} else if val == "\"" {
				break
			} else {
				out += val
			}
		}
		ret := new(Token)
		ret.value = out
		ret.tok_type = "string"
		curr_tok = ret
		return ret
	}

	skip_whitespace := func() {
		//c := input.peek()
		for !input.eof() && strings.Index(" \n\t\r", input.peek()) > -1 {
			input.next()
			//fmt.Println("SKIPPING WS: ", input.peek())
		}
	}

	parse_array := func(c string) *Token { input.next(); return nil }

	parse_object := func(c string) *Token { input.next(); return nil }

	//parse_number := func(c string) *Token { input.next(); return nil }

	parse_special := func(c string) *Token {
		for !input.eof() {
			matched, err := regexp.MatchString("[a-zA-Z_0-9]", input.peek())
			if err != nil || !matched {
				break
			}
			c += input.next()
		}
		ret := new(Token)
		ret.tok_type = "special"
		ret.value = c
		return ret
	}

	parse_next := func() *Token {
		skip_whitespace()
		if input.eof() {
			return nil
		}
		//fmt.Println("BEFORE NEXT: ", input.peek())
		c := input.next()
		//fmt.Println("PARSING NEXT: ", c)
		// check what we should do
		if strings.Index("{}:,", c) > -1 {
			ret := new(Token)
			ret.tok_type = "punctuation"
			ret.value = c
			curr_tok = ret
			return ret
		} else if c == "\"" {
			return parse_string(c)
		} else if c == "[" {
			return parse_array(c)
		} else if c == "{" {
			return parse_object(c)
		} else if strings.Index("tfn", c) > -1 {
			return parse_special(c)
		}
		return nil
	}

	next := func() *Token {
		return parse_next()
	}

	peek := func() *Token {
		if curr_tok == nil {
			curr_tok = parse_next()
			return curr_tok
		}
		return curr_tok
	}
	eof := func() bool { return input.eof() }
	suffocate := func(msg string, list ...string) { input.suffocate(msg) }

	// here go the predicates
	//is_identifier := func(c string) { return regexp.MatchString("", s) }

	ret.next = next
	ret.peek = peek
	ret.eof = eof
	ret.suffocate = suffocate
	return ret
}

func deserialize(json string) (m map[string]interface{}) {
	ret := make(map[string]interface{})
	//ret["hello"] = make(map[string]interface{})
	c := make(chan string)
	//done := make(chan bool)
	spl := strings.Split(json, "")
	fmt.Println("SPLIT: ", spl)
	input := make_iterator(json)

	go func() {
		for !input.eof() {
			val := input.next()
			if val == "\"" {
				fmt.Println("VAL is bracket: ", val)
			}
			fmt.Println("SENT: ", val)
			c <- val
		}
		close(c)
	}()

	// always send one character at a time
	for msg := range c {
		fmt.Println("GOT: ", msg)
	}

	//<-done
	//time.Sleep(time.Second * 10)
	return ret
}

func query(s string, m map[string]interface{}) {
	spl := strings.Split(s, ".")
	curr := m
	for _, val := range spl {
		if curr[val] == nil {
			curr[val] = make(map[string]interface{})
		}
		fmt.Println("Curr now: ", curr)
		curr = curr[val].(map[string]interface{})
	}
	fmt.Println(spl)
}

func main() {
	dat, err := ioutil.ReadFile("./test.json")
	if err != nil {
		fmt.Println("ERR: ", err)
	}
	json := string(dat)
	t := tokenizer(make_iterator(json))
	//a := new([3]*Token)

	for !t.eof() {
		fmt.Println("t: ", t.next())
	}
	//t.next(), t.next(), t.next())

	//m := deserialize(json)
	//fmt.Println("M: ", m, json)
}
