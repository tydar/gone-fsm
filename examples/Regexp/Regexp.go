package main

import (
	"fmt"

	fsm "github.com/tydar/gonefsm"
)

func onlyOne(s string) (*fsm.FSM, error) {
	// onlyOne returns an FSM that matches a regex like ^a?$
	if len(s) != 1 {
		return &fsm.FSM{}, fmt.Errorf("onlyOne must be passed a string of length 1")
	}

	initial := "start"
	final := s + "_end"
	event := fsm.Event{
		From:  "start",
		Input: s,
	}
	events := make(map[fsm.Event]string)
	events[event] = final

	return fsm.NewFSM(initial, events, []string{final}), nil
}

func oneOrMore(s string) (*fsm.FSM, error) {
	// oneOrMore returns an FSM that matches a regex like ^a+?
	if len(s) != 1 {
		return &fsm.FSM{}, fmt.Errorf("oneOrMore must be passed a string of length 1")
	}

	initial := "start"
	accept := s
	event := fsm.Event{
		From:  "start",
		Input: s,
	}
	eventAccept := fsm.Event{
		From:  s,
		Input: s,
	}
	events := make(map[fsm.Event]string)
	events[event] = accept
	events[eventAccept] = accept

	return fsm.NewFSM(initial, events, []string{accept}), nil
}

func anyOfOneOrMore(s string) (*fsm.FSM, error) {
	// anyOfOneOrMore returns an FSM that matches a regex like ^[s]+?
	// e.g. one or more of any of the characters
	// does not support [A-Z] syntax
	if len(s) == 0 {
		return &fsm.FSM{}, fmt.Errorf("anyOfOneOrMore requires at least one character")
	}

	initial := "start"
	events := make(map[fsm.Event]string)
	var accept []string

	// for every character in the string, we want an event from initial -> that char accept state
	// and an event from that char accept state to itself
	for _, c := range s {
		initToSelf := fsm.Event{
			Input: string(c),
			From:  initial,
		}
		selfToSelf := fsm.Event{
			Input: string(c),
			From:  string(c),
		}

		// check for an existing value prior to storing a new one
		// actually unnecessary since it would just enter the same value
		_, ok := events[initToSelf]
		if !ok {
			events[initToSelf] = string(c)
		}

		_, ok = events[selfToSelf]
		if !ok {
			events[selfToSelf] = string(c)
		}
		accept = append(accept, string(c))
	}

	return fsm.NewFSM(initial, events, accept), nil
}

func processString(f *fsm.FSM, s string) error {
	// processString feeds a string one rune at a time to an FSM
	// will return an error if a character does not match
	// will return nil if all characters accepted. then have to check FSM state to determine acceptance
	for _, char := range s {
		err := f.Event(string(char))
		if err != nil {
			// if we get an error, that means it did not match
			// so just stop and return the error
			return err
		}
	}
	return nil
}

func checkMatch(f *fsm.FSM, s string) bool {
	// checkMatch returns true if the given FSM matches the string
	// or false if it does not (e.g. f.Accepted != true or get an error from processString)
	err := processString(f, s)
	if err == nil {
		return f.Accepted()
	}
	return false
}

func main() {
	f1, err1 := onlyOne("a")
	if err1 != nil {
		panic(err1)
	}

	f2, err2 := oneOrMore("a")
	if err2 != nil {
		panic(err2)
	}

	f3, err3 := anyOfOneOrMore("abc")
	if err3 != nil {
		panic(err3)
	}

	good := "a"
	bad := "aa"
	worse := "bbbb"
	mixed := "aabb"

	fmt.Printf("The string good '%s' was accepted by the FSM f1? %v\n", good, checkMatch(f1, good))

	f1.Reset()

	fmt.Printf("The string bad '%s' was accepted by the FSM f1? %v\n", bad, checkMatch(f1, bad))

	fmt.Printf("The string bad '%s' was accepted by the FSM f2? %v\n", bad, checkMatch(f2, bad))

	f2.Reset()

	fmt.Printf("The string worse '%s' was accepted by the FSM f2? %v\n", worse, checkMatch(f2, worse))

	fmt.Printf("The string worse '%s' was accepted by the FSM f3? %v\n", worse, checkMatch(f3, worse))

	f3.Reset()
	fmt.Printf("The string mixed '%s' was accepted by the FSM f3? %v\n", mixed, checkMatch(f3, mixed))
}
