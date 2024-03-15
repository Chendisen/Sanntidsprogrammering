package door_open_timer

import (
	"Sanntid/elevator"
	"Sanntid/timer"
	"Sanntid/world_view"
	"Sanntid/fsm"
)

func CheckDoorOpenTimeout(elev *elevator.Elevator, worldView *world_view.WorldView, myIP string, tmr *timer.Timer, watchdog *timer.Timer) {
	for {
		if elev.DoorObstructed {
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
		}
		if tmr.Timer_timedOut(elev.Config.DoorOpenDuration_s) {
			tmr.Timer_stop()
			fsm.Fsm_onDoorTimeout(elev, worldView, myIP, tmr, watchdog)
		}
	}
}