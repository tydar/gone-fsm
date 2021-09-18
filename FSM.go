package FSM

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
}

func (f *FSM) GetEvent(start string, input string) (Event, error) {
	// GetEvent returns a single transition and an error value based on a starting state and an input symbol
	// GetEvent assumes there is only one transition that matches the starting state and the input state
	for k := range f.Events {
		if k.From == start && k.Input == input {
			return k, nil
		}
	}
	return Event{}, fmt.Errorf("no transtion from state %s with input %s found", start, input)
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

func NewFSM(initial string, events map[Event]string) *FSM {
	// NewFSM creates a new FSM by receiving an initial state and a map describing all possible transitions
	states := make(map[string]int)
	for k, v := range events {
		// if destination state not yet added, initialize
		for _, present := states[v]; !present; {
			states[v] = 0
		}
		// if start state not present, initialize
		for _, present := states[k.From]; !present; {
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
	}
}
