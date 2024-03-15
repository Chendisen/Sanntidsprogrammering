package door_open_timer

import (
	// "Sanntid/driver"
	"Sanntid/world_view"
	"Sanntid/elevator"
	"Sanntid/fsm"
	"Sanntid/timer"
)

func CheckDoorOpenTimeout(elev *elevator.Elevator, myIP string, tmr *timer.Timer, watchdog *timer.Timer, upd_request chan world_view.UpdateRequest) {
	for {
		if elev.DoorObstructed {
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
		}
		if tmr.Timer_timedOut(elev.Config.DoorOpenDuration_s) {
			tmr.Timer_stop()
			fsm.Fsm_onDoorTimeout(elev, myIP, tmr, watchdog, upd_request)
		}
	}
}