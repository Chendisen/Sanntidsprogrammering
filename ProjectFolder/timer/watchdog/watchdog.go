package watchdog

import (
	"Sanntid/elevator"
	"Sanntid/timer"
)

func CheckWatchdogTimeout(tmr *timer.Timer, elevState *elevator.Elevator, dead chan<- bool) {
	
	for {
		if tmr.Timer_timedOut(timer.WATCHDOG_TimeoutTime) {
			if elevState.Behaviour == elevator.EB_Moving && !elevState.DoorObstructed{
				tmr.Timer_stop()
				dead <- true
			} else {
				tmr.Timer_start(timer.WATCHDOG_TimeoutTime)
			}
		}
	}
}


