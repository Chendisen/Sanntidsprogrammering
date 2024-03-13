package watchdog

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/timer"
)

const watchdogTime float64 = 10

func Watchdog(tmr *timer.Timer, es *elevator.Elevator, dead chan<- bool) {
	timeOut := make(chan bool)
	go tmr.TimeOut(timeOut)

	for range timeOut {
		if es.Dirn != driver.MD_Stop {
			dead <- true
		} else {
			timer.Timer_start(tmr, watchdogTime)
		}
	}
}


