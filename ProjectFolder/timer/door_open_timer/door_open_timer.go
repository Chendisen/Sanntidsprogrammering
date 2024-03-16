package door_open_timer

import (
	"Sanntid/elevator"
	. "Sanntid/resources/update_request"
	"Sanntid/timer"
)

func CheckDoorOpenTimeout(elev *elevator.Elevator, myIP string, tmr *timer.Timer, watchdog *timer.Timer, upd_request chan UpdateRequest) {
	for {
		if elev.DoorObstructed {
			tmr.TimerStart(elev.Config.DoorOpenDuration_s)
		}
		if tmr.TimerTimedOut(elev.Config.DoorOpenDuration_s) {
			tmr.TimerStop()
			elevator.Fsm_onDoorTimeout(elev, myIP, tmr, watchdog, upd_request)
		}
	}
}
