package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
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
// - punctuation : or ,
// - string
// - number
// - object - basically just the curly brackets
// - array  - basically just the squared brackets
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

	parse_array := func(c string) *Token {
		ret := new(Token)
		ret.tok_type = "array"
		ret.value = c
		return ret
	}

	parse_object := func(c string) *Token {
		ret := new(Token)
		if c == "{" {
			ret.tok_type = "object_o"
		} else {
			ret.tok_type = "object_c"
		}
		ret.value = c
		return ret
	}

	parse_number := func(c string) *Token {
		// first parse integer value
		for !input.eof() {
			matched, err := regexp.MatchString("[0-9\\.]", input.peek())
			if err != nil || !matched {
				break
			}
			c += input.next()
		}
		// parse part after dot
		if input.peek() == "." {
			c += input.next()
			for !input.eof() {
				matched, err := regexp.MatchString("[0-9]", input.peek())
				if err != nil || !matched {
					break
				}
				c += input.next()
			}
		}
		// parse optional exponential part
		if input.peek() == "e" || input.peek() == "E" {
			c += input.next()
			// TODO: add errors
			if input.peek() == "+" || input.peek() == "-" {
				c += input.next()
				for !input.eof() {
					matched, err := regexp.MatchString("[0-9]", input.peek())
					if err != nil || !matched {
						break
					}
					c += input.next()
				}
			}
		}
		ret := new(Token)
		ret.tok_type = "number"
		ret.value = c
		return ret
	}

	parse_special := func(c string) *Token {
		for !input.eof() {
			matched, err := regexp.MatchString("[a-z]", input.peek())
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

	parse_punctuation := func(c string) *Token {
		ret := new(Token)
		if c == ":" {
			ret.tok_type = "colon"
		} else {
			ret.tok_type = "comma"
		}
		ret.value = c
		curr_tok = ret
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
		if strings.Index(":,", c) > -1 {
			return parse_punctuation(c)
		} else if c == "\"" {
			return parse_string(c)
		} else if strings.Index("[]", c) > -1 {
			return parse_array(c)
		} else if strings.Index("{}", c) > -1 {
			return parse_object(c)
		} else if strings.Index("tfn", c) > -1 {
			return parse_special(c)
		} else if strings.Index("0123456789-", c) > -1 {
			return parse_number(c)
		}
		return nil
	}

	next := func() *Token {
		curr_tok = parse_next()
		return curr_tok
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

type StateMachine struct {
	states  map[string]*State
	actions map[string]map[string]func(t *Token)
	yield   func() map[string]interface{}
}

type State struct {
	name        string
	transitions map[string]string
}

type StateConfig struct {
	name  string
	conns []string
}

type Transition struct {
	target string
	exec   func(tok *Token)
}

func add_state(name string, trans []string) *State {
	s := new(State)
	s.name = name
	s.transitions = make(map[string]string)
	for _, t := range trans {
		spl := strings.Split(t, " -> ")
		s.transitions[spl[0]] = spl[1]
	}
	return s
}

func add_state1(m map[string]*State, sc *StateConfig) {
	s := new(State)
	s.name = sc.name
	s.transitions = make(map[string]string)
	for _, t := range sc.conns {
		spl := strings.Split(t, " -> ")
		s.transitions[spl[0]] = spl[1]
	}
	m[sc.name] = s
}

func add_action(m map[string]map[string]func(t *Token), from_state, to_state string, action func(t *Token)) {
	m[from_state] = make(map[string]func(t *Token))
	m[from_state][to_state] = action
}

func new_machine() *StateMachine {
	state_machine := new(StateMachine)
	state_machine.states = make(map[string]*State)

	// define the states and transitions
	sc := []*StateConfig{&StateConfig{"initial", []string{"object_o -> await_key_end"}},
		&StateConfig{"await_key_end", []string{"object_c -> end?", "string -> await_colon"}},
		&StateConfig{"await_colon", []string{"colon -> await_expr"}},
		&StateConfig{"await_expr", []string{"string -> comma_end?", "number -> comma_end?", "special -> comma_end?"}},
		&StateConfig{"comma_end?", []string{"object_c -> end?", "comma -> new_pair"}},
		&StateConfig{"new_pair", []string{"string -> await_colon"}},
		&StateConfig{"end?", []string{}}}

	for _, conf := range sc {
		add_state1(state_machine.states, conf)
	}

	// ==========================================
	// Actions
	// ==========================================
	var ret map[string]interface{}
	var curr map[string]interface{}
	var key string
	//var value interface{}

	setup := func(t *Token) {
		ret = make(map[string]interface{})
		curr = ret
		//fmt.Println("DONE WITH SETUP: ", ret)
	}

	set_key := func(t *Token) {
		key = t.value
		//fmt.Println("SET KEY TO: ", key)
	}

	add_entry := func(t *Token) {
		var val interface{}
		if t.tok_type == "special" {
			if t.value == "true" {
				val = true
			} else if t.value == "false" {
				val = false
			} else if t.value == "null" {
				val = nil
			}
		} else if t.tok_type == "number" {
			// parse either float or integer
			m1, _ := regexp.MatchString("[0-9]+\\.[0-9]+", t.value)
			m2, _ := regexp.MatchString("[Ee][\\-\\+][0-9]+$", t.value)
			if m1 || m2 {
				val, _ = strconv.ParseFloat(t.value, 64)
			} else {
				val, _ = strconv.Atoi(t.value)
			}
		} else {
			val = t.value
		}
		curr[key] = val
		key = ""
	}

	//state_machine.actions = new(Action)
	state_machine.actions = make(map[string]map[string]func(t *Token))

	add_action(state_machine.actions, "initial", "await_key_end", setup)
	add_action(state_machine.actions, "await_key_end", "await_colon", set_key)
	add_action(state_machine.actions, "await_expr", "comma_end?", add_entry)
	add_action(state_machine.actions, "new_pair", "await_colon", set_key)

	yield := func() map[string]interface{} {
		return ret
	}

	state_machine.yield = yield

	return state_machine
}

func json_mapper(tokens *Tokenizer) map[string]interface{} {
	/*var ret map[string]interface{}
	var curr map[string]interface{}
	var key string
	var value interface{}*/

	// functions that we'll need
	/* - setup -> initializes ret and curr
	 * - set_key -> sets the key variable
	 * - add_entry -> adds a parsed entry to curr and then flushes key
	 * - close_object -> closes the current object and pops from stack
	 * - wait_for_more -> waits for more expressions (comma)
	 */

	// state map
	curr_state := "initial"
	prev_state := ""
	machine := new_machine()

	var t *Token
	// loop through the tokens
	for !tokens.eof() {
		t = tokens.next()
		// now dispatch the mapper
		fmt.Println("STATE NOW => ", curr_state)
		if tmp, found := machine.states[curr_state]; found && t != nil {
			if s, f := tmp.transitions[t.tok_type]; f {
				prev_state = curr_state
				curr_state = s
				// now check whether action exists
				if a, found := machine.actions[prev_state]; found {
					if x, f := a[curr_state]; f {
						x(t)
					}
				}
			}
		}

		fmt.Println("DISPATCHING: ", t)
	}

	return machine.yield()
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

func query(s string, m map[string]interface{}) bool {
	spl := strings.Split(s, ".")
	curr := m
	for _, val := range spl {
		if curr[val] == nil {
			return false
		}
		//fmt.Println("Curr now: ", curr)
		curr = curr[val].(map[string]interface{})
	}
	//fmt.Println(spl)
	return true
}

func get_state(name, to string, m map[string]*State) (string, bool) {
	s, has := m[name].transitions[to]
	fmt.Println("S => ", s)
	return s, has
}

func main() {
	dat, err := ioutil.ReadFile("./test.json")
	if err != nil {
		fmt.Println("ERR: ", err)
	}
	json := string(dat)
	t := tokenizer(make_iterator(json))
	//a := new([3]*Token)
	parsed := json_mapper(t)

	fmt.Println("TEST PARSED NUM: ", parsed["num"].(float64)+10, "\n", parsed)

	/*for !t.eof() {
		fmt.Println("t: ", t.peek())
		t.next()
	}*/
	//t.next(), t.next(), t.next())

	//m := deserialize(json)
	//fmt.Println("M: ", m, json)
}
