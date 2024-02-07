package requests

import (
	"Sanntid/elevator"
)

type DirnBehaviourPair struct {
	dirn      elevator_io.Dirn
	behaviour elevator.ElevatorBehaviour
}

func requests_above(e elevator.Elevator) int {
	for floor := e.Floor + 1; floor < elevator_io.N_FLOORS; floor++ {
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			if e.Request[floor][btn] {
				return 1
			}
		}
	}
	return 0
}

func requests_below(e elevator.Elevator) int {
	for floor := 0; floor < e.Floor; floor++ {
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			if e.Request[floor][btn] {
				return 1
			}
		}
	}
	return 0
}

func requests_here(e elevator.Elevator) int {
	for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
		if e.Request[e.Floor][btn] {
			return 1
		}
	}
	return 0
}

func Requests_chooseDirection(e elevator.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case elevator_io.D_Up:
		if requests_above(e) {
			return DirnBehaviourPair({elevator_io.D_Up, elevator.EB_Moving})
		}
		if requests_here(e) {
			return DirnBehaviourPair({elevator_io.D_Down, elevator.EB_DoorOpen})
		}
		if requests_below(e) {
			return DirnBehaviourPair({elevator_io.D_Down, elevator.EB_Moving})
		}
		return DirnBehaviourPair({elevator_io.D_Stop, elevator.EB_Idle})
	case elevator_io.D_Down:
		if requests_below(e) {
			return DirnBehaviourPair({elevator_io.D_Down, elevator.EB_Moving})
		}
		if requests_here(e) {
			return DirnBehaviourPair({elevator_io.D_Up, elevator.EB_DoorOpen})
		}
		if requests_above(e) {
			return DirnBehaviourPair({elevator_io.D_Up, elevator.EB_Moving})
		}
		return DirnBehaviourPair({elevator_io.D_Stop, elevator.EB_Idle})
	case elevator_io.D_Stop: // There should only be one request in the Stop case. Checking up or down first is arbitrary
		if requests_here(e) {
			return DirnBehaviourPair({elevator_io.D_Stop, elevator.EB_DoorOpen})
		}
		if requests_above(e) {
			return DirnBehaviourPair({elevator_io.D_Up, elevator.EB_Moving})
		}
		if requests_below(e) {
			return DirnBehaviourPair({elevator_io.D_Down, elevator.EB_Moving})
		}
		return DirnBehaviourPair({elevator_io.D_Stop, elevator.EB_Idle})
	default:
		return DirnBehaviourPair({elevator_io.D_Stop, elevator.EB_Idle})
	}
}

func Requests_shouldStop(e elevator.Elevator) int {
	switch e.Dirn {
	case elevator_io.D_Down:
		return 
			((e.Requests[e.Floor][elevator_io.B_HallDown]) || 
			(e.Requests[e.Floor][elevator_io.B_Cab]) || 
			!requests_below(e))
	case elevator_io.D_Up:
		return 
			((e.Requests[e.Floor][elevator_io.B_HallUp]) || 
			(e.Requests[e.Floor][elevator_io.B_Cab]) || 
			!requests_above(e))
	case elevator_io.D_Stop:
		return 1
	default:
		return 1
	}
}

func Requests_shouldClearImmediately(e elevator.Elevator, btn_floor int, btn_type elevator_io.Button) int {
	switch e.Config.ClearRequestVariant{
	case CV_all:
		return e.Floor == btn_floor
	case CV_InDirn:
		return e.Floor == btn_floor &&
				(
					(e.Dirn == elevator_io.D_Up && btn_type == elevator_io.B_HallUp) ||
					(e.Dirn == elevator_io.D_Down && btn_type == elevator_io.B_HallDown) ||
					(e.Dirn == elevator_io.D_Stop) ||
					(btn_type == elevator_io.B_Cab)
				)
	default:
		return 0	
	}
}

func Requests_clearAtCurrentFloor(e elevator.Elevator) elevator.Elevator {
	switch e.Config.ClearRequestVariant{
	case CV_all:
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			e.Request[e.Floor][btn] = 0
		}
	case CV_InDirn:
		e.Request[e.Floor][elevator_io.B_Cab] = 0
		switch e.Dirn {
		case elevator_io.D_Up:
			if !requests_above(e) && !e.Requests[e.Floor][elevator_io.B_HallUp] {
				e.Requests[e.Floor][elevator_io.B_HallDown] = 0
			}
			e.Requests[e.Floor][elevator_io.B_HallUp] = 0
		case elevator_io.D_Down:
			if !requests_below(e) && !e.Requests[e.Floor][elevator_io.B_HallDown] {
				e.Requests[e.Floor][elevator_io.B_HallUp] = 0
			}
			e.Requests[e.Floor][elevator_io.B_Halldown] = 0
		case elevator_io.D_Stop:
			e.Requests[e.Floor][elevator_io.B_HallUp] = 0
			e.Requests[e.Floor][elevator_io.B_HallDown] = 0
		default:
			e.Requests[e.Floor][elevator_io.B_HallUp] = 0
			e.Requests[e.Floor][elevator_io.B_HallDown] = 0
		}
	}
	return e
}