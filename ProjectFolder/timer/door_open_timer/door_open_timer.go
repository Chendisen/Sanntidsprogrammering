package door_open_timer

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/fsm"
	"Sanntid/timer"
)

func CheckDoorOpenTimeout(elev *elevator.Elevator, myIP string, tmr *timer.Timer, watchdog *timer.Timer, set_behaviour chan<- elevator.ElevatorBehaviour, set_direction chan<- driver.MotorDirection, finished_request_at_floor chan<- driver.ButtonEvent) {
	for {
		if elev.DoorObstructed {
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
		}
		if tmr.Timer_timedOut(elev.Config.DoorOpenDuration_s) {
			tmr.Timer_stop()
			fsm.Fsm_onDoorTimeout(elev, myIP, tmr, watchdog, set_behaviour, set_direction, finished_request_at_floor)
		}
	}
}