package cyclic_counter

const MAX int = 15

type Counter struct {
	Value int
	Max int
	ShouldReset bool
}

func MakeCounter(max int) Counter{
	return Counter{Value: 0, Max: max, ShouldReset: false}
}

func Reset(counter *Counter){
	counter.Value = 0
	counter.ShouldReset = false
}

func ShouldUpdate(recieved Counter, current Counter) bool{
	var shouldUpdate bool = false

	if(recieved.Value > current.Value) {
		shouldUpdate = true
	} else if(recieved.Value < current.Value && current.ShouldReset){
		shouldUpdate = true
	} 

	return shouldUpdate
}

func Increment(counter *Counter){
	if(counter.Value == counter.Max){
		Reset(counter)
	} else {
		counter.Value++
	}
}

func UpdateValue(counter *Counter, value int){
	if(counter.Value > value){
		counter.Value = value
		counter.ShouldReset = false
	} else {
		counter.Value = value
	}
}

func (cc Counter) ToBool() bool{
	return cc.Value%2 == 1
}