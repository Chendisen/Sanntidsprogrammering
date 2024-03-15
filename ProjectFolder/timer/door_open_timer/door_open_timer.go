package door_open_timer

import (
	// "Sanntid/driver"
	"Sanntid/elevator"
	. "Sanntid/resources/update_request"
	"Sanntid/timer"
)

func CheckDoorOpenTimeout(elev *elevator.Elevator, myIP string, tmr *timer.Timer, watchdog *timer.Timer, upd_request chan UpdateRequest) {
	for {
		//fmt.Printf("Timeout time: %f Current time %f\n", tmr.TimerEndTime, timer.Get_current_time())
		if elev.DoorObstructed {
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
		}
		if tmr.Timer_timedOut(elev.Config.DoorOpenDuration_s) {
			tmr.Timer_stop()
			elevator.Fsm_onDoorTimeout(elev, myIP, tmr, watchdog, upd_request)
		}
	}
}
