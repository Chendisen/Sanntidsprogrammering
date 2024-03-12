package cyclic_counter

import "fmt"

const MAX int = 3

type Counter struct {
	Value       int
	Max         int
	ShouldReset bool
}

func MakeCounter(max int) Counter {
	return Counter{Value: 0, Max: max, ShouldReset: false}
}

func Reset(counter *Counter) {
	counter.Value = 0
	counter.ShouldReset = false
}

func ShouldUpdate(recieved Counter, current Counter) bool {
	var shouldUpdate bool = false

	if recieved.Value > current.Value {
		shouldUpdate = true
	} else if recieved.Value < current.Value && current.ShouldReset {
		shouldUpdate = true
	}

	return shouldUpdate
}

func Increment(counter *Counter) {
	if counter.Value == counter.Max {
		counter.PrintCounter()
		Reset(counter)
	} else {
		counter.Value++
	}
	if counter.Value == counter.Max {
		counter.ShouldReset = true
	}
}

func UpdateValue(counter *Counter, value int) {
	if counter.Value > value {
		counter.Value = value
		counter.ShouldReset = false
	} else {
		counter.Value = value
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
