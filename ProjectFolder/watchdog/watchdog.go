package watchdog

import (
	"Sanntid/elevator"
	"Sanntid/timer"
)

const watchdogTime float64 = 5

func Watchdog(tmr *timer.Timer, es *elevator.Elevator, dead chan<- bool) {
	
	for {
		if tmr.Timer_timedOut(watchdogTime) {
			if es.Behaviour == elevator.EB_Moving && !es.DoorObstructed{
				tmr.Timer_stop()
				dead <- true
			} else {
				tmr.Timer_start(watchdogTime)
			}
		}
	}
}


