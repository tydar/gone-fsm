package main

import (
	"fmt"

	fsm "github.com/tydar/gonefsm"
)

func main() {
	events := make(map[fsm.Event](string))
	onFlip := fsm.Event{
		Input: "flip",
		From:  "on",
	}
	events[onFlip] = "off"
	offFlip := fsm.Event{
		Input: "flip",
		From:  "off",
	}
	events[offFlip] = "on"
	f := fsm.NewFSM("off", events)
	fmt.Printf("Current state: %s\n", f.CurrentState)
	fmt.Println("Flipping the switch...")
	err := f.Event("flip")

	if err != nil {
		panic(err)
	}

	fmt.Printf("Current state: %s\n", f.CurrentState)
}
