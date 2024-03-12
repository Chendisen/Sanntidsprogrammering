package cyclic_counter

import (
	"fmt"
)

const MAX int = 3

type Counter struct {
	Value       int
	Max         int
	ShouldReset bool
}

func MakeCounter(max int) Counter {
	return Counter{Value: 0, Max: max, ShouldReset: false}
}

func (cc *Counter) Reset() {
	cc.Value = 0
	cc.ShouldReset = false
}

func (current *Counter) ShouldUpdate(recieved Counter) bool {
	var shouldUpdate bool = false

	if recieved.Value > current.Value {
		shouldUpdate = true
	} else if recieved.Value < current.Value && current.ShouldReset {
		shouldUpdate = true
	}

	return shouldUpdate
}

func (cc *Counter) Increment() {
	if cc.Value == cc.Max {
		cc.PrintCounter()
		cc.Reset()
	} else {
		cc.Value++
	}
	if cc.Value == cc.Max {
		cc.ShouldReset = true
	}
}

func (cc *Counter) UpdateValue(value int) {
	if cc.Value > value {
		cc.Value = value
		cc.ShouldReset = false
	} else {
		cc.Value = value
	}
}

func (cc Counter) ToBool() bool {
	return cc.Value == 2 || cc.Value == 3
}

func (cc Counter) PrintCounter() {
	fmt.Println("Counter state: ")
	fmt.Printf("	Value: %d\n", cc.Value)
	fmt.Printf("	Reset: %t\n", cc.ShouldReset)
}
