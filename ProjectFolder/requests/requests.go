package requests

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/elevator_io"
)

type DirnBehaviourPair struct {
	dirn      driver.MotorDirection
	behaviour elevator.ElevatorBehaviour
}

func requests_above(e elevator.Elevator) int {
	for floor := e.Floor + 1; floor < elevator_io.N_FLOORS; floor++ {
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			if driver.IntToBool(e.Request[floor][btn]) {
				return 1
			}
		}
	}
	return 0
}

func requests_below(e elevator.Elevator) int {
	for floor := 0; floor < e.Floor; floor++ {
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			if driver.IntToBool(e.Request[floor][btn]) {
				return 1
			}
		}
	}
	return 0
}

func requests_here(e elevator.Elevator) int {
	for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
		if driver.IntToBool(e.Request[e.Floor][btn]) {
			return 1
		}
	}
	return 0
}

func Requests_chooseDirection(e elevator.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case driver.MD_Up:
		if driver.IntToBool(requests_above(e)) {
			return DirnBehaviourPair{driver.MD_Up, elevator.EB_Moving}
		}
		if driver.IntToBool(requests_here(e)) {
			return DirnBehaviourPair{driver.MD_Down, elevator.EB_DoorOpen}
		}
		if driver.IntToBool(requests_below(e)) {
			return DirnBehaviourPair{driver.MD_Down, elevator.EB_Moving}
		}
		return DirnBehaviourPair{driver.MD_Stop, elevator.EB_Idle}
	case driver.MD_Down:
		if driver.IntToBool(requests_below(e)) {
			return DirnBehaviourPair{driver.MD_Down, elevator.EB_Moving}
		}
		if driver.IntToBool(requests_here(e)) {
			return DirnBehaviourPair{driver.MD_Up, elevator.EB_DoorOpen}
		}
		if driver.IntToBool(requests_above(e)) {
			return DirnBehaviourPair{driver.MD_Up, elevator.EB_Moving}
		}
		return DirnBehaviourPair{driver.MD_Stop, elevator.EB_Idle}
	case driver.MD_Stop: // There should only be one request in the Stop case. Checking up or down first is arbitrary
		if driver.IntToBool(requests_here(e)) {
			return DirnBehaviourPair{driver.MD_Stop, elevator.EB_DoorOpen}
		}
		if driver.IntToBool(requests_above(e)) {
			return DirnBehaviourPair{driver.MD_Up, elevator.EB_Moving}
		}
		if driver.IntToBool(requests_below(e)) {
			return DirnBehaviourPair{driver.MD_Down, elevator.EB_Moving}
		}
		return DirnBehaviourPair{driver.MD_Stop, elevator.EB_Idle}
	default:
		return DirnBehaviourPair{driver.MD_Stop, elevator.EB_Idle}
	}
}

func Requests_shouldStop(e elevator.Elevator) int {
	switch e.Dirn {
	case driver.MD_Down:
		return driver.BoolToInt((driver.IntToBool(e.Request[e.Floor][driver.BT_HallDown])) ||
								(driver.IntToBool(e.Request[e.Floor][driver.BT_Cab])) ||
								!driver.IntToBool(requests_below(e)))
	case driver.MD_Up:
		return driver.BoolToInt((driver.IntToBool(e.Request[e.Floor][driver.BT_HallUp])) ||
								(driver.IntToBool(e.Request[e.Floor][driver.BT_Cab])) ||
								!driver.IntToBool(requests_above(e)))
	case driver.MD_Stop:
		return 1
	default:
		return 1
	}
}

func Requests_shouldClearImmediately(e elevator.Elevator, btn_floor int, btn_type elevator_io.Button) int {
	switch e.Config.ClearRequestVariant {
	case elevator.CV_all:
		return driver.BoolToInt(e.Floor == btn_floor)
	case elevator.CV_InDirn:
		return e.Floor == btn_floor &&
			((e.Dirn == driver.MD_Up && btn_type == driver.BT_HallUp) ||
				(e.Dirn == driver.MD_Down && btn_type == driver.BT_HallDown) ||
				(e.Dirn == driver.MD_Stop) ||
				(btn_type == driver.BT_Cab))
	default:
		return 0
	}
}

func Requests_clearAtCurrentFloor(e elevator.Elevator) elevator.Elevator {
	switch e.Config.ClearRequestVariant {
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
