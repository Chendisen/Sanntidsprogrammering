package requests

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/world_view"
)

type DirnBehaviourPair struct {
	Dirn      driver.MotorDirection
	Behaviour elevator.ElevatorBehaviour
}

func requests_above(e elevator.Elevator) int {
	for floor := e.Floor + 1; floor < driver.N_FLOORS; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			if driver.IntToBool(e.Request[floor][btn]) {
				return 1
			}
		}
	}
	return 0
}

func requests_below(e elevator.Elevator) int {
	for floor := 0; floor < e.Floor; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			if driver.IntToBool(e.Request[floor][btn]) {
				return 1
			}
		}
	}
	return 0
}

func requests_here(e elevator.Elevator) int {
	for btn := 0; btn < driver.N_BUTTONS; btn++ {
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

func Requests_shouldStop(e elevator.Elevator) bool {
	switch e.Dirn {
	case driver.MD_Down:
		return (driver.IntToBool(e.Request[e.Floor][driver.BT_HallDown])) ||
			(driver.IntToBool(e.Request[e.Floor][driver.BT_Cab])) ||
			!driver.IntToBool(requests_below(e))
	case driver.MD_Up:
		return (driver.IntToBool(e.Request[e.Floor][driver.BT_HallUp])) ||
			(driver.IntToBool(e.Request[e.Floor][driver.BT_Cab])) ||
			!driver.IntToBool(requests_above(e))
	case driver.MD_Stop:
		return true
	default:
		return true
	}
}

func Requests_shouldClearImmediately(e elevator.Elevator, btn_floor int, btn_type driver.ButtonType) bool {
	switch e.Config.ClearRequestVariant {
	case elevator.CV_all:
		return e.Floor == btn_floor
	case elevator.CV_InDirn:
		return (e.Floor == btn_floor &&
			((e.Dirn == driver.MD_Up && btn_type == driver.BT_HallUp) ||
				(e.Dirn == driver.MD_Down && btn_type == driver.BT_HallDown) ||
				(e.Dirn == driver.MD_Stop) ||
				(btn_type == driver.BT_Cab)))
	default:
		return false
	}
}

func Requests_clearAtCurrentFloor(e *elevator.Elevator, wld_view *world_view.WorldView, myIP string) {
	switch e.Config.ClearRequestVariant {
	case elevator.CV_all:
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			e.Request[e.Floor][btn] = 0
			wld_view.ClearRequestAtFloor(myIP, e.Floor, btn)
		}
	case elevator.CV_InDirn:
		e.Request[e.Floor][driver.BT_Cab] = 0
		wld_view.ClearRequestAtFloor(myIP, e.Floor, driver.BT_Cab)
		switch e.Dirn {
		case driver.MD_Up:
			if !driver.IntToBool(requests_above(*e)) && !driver.IntToBool(e.Request[e.Floor][driver.BT_HallUp]) {
				e.Request[e.Floor][driver.BT_HallDown] = 0
				wld_view.ClearRequestAtFloor(myIP, e.Floor, driver.BT_HallDown)
			}
			e.Request[e.Floor][driver.BT_HallUp] = 0
			wld_view.ClearRequestAtFloor(myIP, e.Floor, driver.BT_HallUp)
		case driver.MD_Down:
			if !driver.IntToBool(requests_below(*e)) && !driver.IntToBool(e.Request[e.Floor][driver.BT_HallDown]) {
				e.Request[e.Floor][driver.BT_HallUp] = 0
				wld_view.ClearRequestAtFloor(myIP, e.Floor, driver.BT_HallUp)
			}
			e.Request[e.Floor][driver.BT_HallDown] = 0
			wld_view.ClearRequestAtFloor(myIP, e.Floor, driver.BT_HallDown)
		case driver.MD_Stop:
			e.Request[e.Floor][driver.BT_HallUp] = 0
			e.Request[e.Floor][driver.BT_HallDown] = 0
			wld_view.ClearRequestAtFloor(myIP, e.Floor, driver.BT_HallUp)
			wld_view.ClearRequestAtFloor(myIP, e.Floor, driver.BT_HallDown)
		default:
			e.Request[e.Floor][driver.BT_HallUp] = 0
			e.Request[e.Floor][driver.BT_HallDown] = 0
			wld_view.ClearRequestAtFloor(myIP, e.Floor, driver.BT_HallUp)
			wld_view.ClearRequestAtFloor(myIP, e.Floor, driver.BT_HallDown)
		}
	}
}
