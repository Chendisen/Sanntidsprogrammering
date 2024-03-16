package watchdog

import (
	"Sanntid/elevator"
	"Sanntid/timer"
)

func CheckWatchdogTimeout(tmr *timer.Timer, elevState *elevator.Elevator, dead chan<- bool) {
	
	for {
		if tmr.TimerTimedOut(timer.WATCHDOG_TimeoutTime) {
			if elevState.Behaviour == elevator.EB_Moving && !elevState.DoorObstructed{
				tmr.TimerStop()
				dead <- true
			} else {
				tmr.TimerStart(timer.WATCHDOG_TimeoutTime)
			}
		}
	}
}


