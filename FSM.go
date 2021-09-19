package fsm

import (
	"fmt"
)

// FSM package implementing a basic set of functions for a deterministic FSM

type Event struct {
	Input string
	From  string
}

type FSM struct {
	States       []string
	Events       map[Event]string
	InitialState string
	CurrentState string
	AcceptStates []string
}

func (f *FSM) GetEvent(start string, input string) (Event, error) {
	// GetEvent returns a single event and an error value based on a starting state and an input symbol
	e := Event{
		From:  start,
		Input: input,
	}
	_, present := f.Events[e]

	if present {
		return e, nil
	}

	return Event{}, &UnmatchedEventError{
		Input:        input,
		CurrentState: f.CurrentState,
	}
}

func (f *FSM) Event(input string) error {
	// Event will transition the state as appropriate for the input value
	// If no transition from the current state with the current input is found an error is returned
	e, err := f.GetEvent(f.CurrentState, input)

	if err != nil {
		return err
	}

	f.CurrentState = f.Events[e]
	return nil
}

func (f *FSM) Accepted() bool {
	// Accepted returns true if the current state of the FSM is one of the accept states
	for i := range f.AcceptStates {
		if f.AcceptStates[i] == f.CurrentState {
			return true
		}
	}
	return false
}

func NewFSM(initial string, events map[Event]string, accept []string) *FSM {
	// NewFSM creates a new FSM by receiving an initial state and a map describing all possible transitions
	states := make(map[string]int)
	for k, v := range events {
		// if destination state not yet added, initialize
		_, present := states[v]
		if !present {
			states[v] = 0
		}
		// if start state not present, initialize
		_, present = states[k.From]
		if !present {
			states[k.From] = 0
		}
		states[k.From] += 1
		states[v] += 1
	}

	// now that we've derived the list of States from the events provided
	// transform the list of keys to a slice
	var statesSlice []string
	for k := range states {
		statesSlice = append(statesSlice, k)
	}

	return &FSM{
		InitialState: initial,
		CurrentState: initial,
		Events:       events,
		States:       statesSlice,
		AcceptStates: accept,
	}
}

// Error definitions begin here
type UnmatchedEventError struct {
	Input        string
	CurrentState string
}

func (e *UnmatchedEventError) Error() string {
	return fmt.Sprintf("no valid transition with input %s from state %s", e.Input, e.CurrentState)
}
