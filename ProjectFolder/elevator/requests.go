package elevator

import (
	"Sanntid/resources/driver"
	. "Sanntid/resources/update_request"
)

type DirnBehaviourPair struct {
	Dirn      driver.MotorDirection
	Behaviour ElevatorBehaviour
}


func requests_above(elev Elevator) int {
	for floor := elev.Floor + 1; floor < driver.N_FLOORS; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			if intToBool(elev.Request[floor][btn]) {
				return 1
			}
		}
	}
	return 0
}

func requests_below(elev Elevator) int {
	for floor := 0; floor < elev.Floor; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			if intToBool(elev.Request[floor][btn]) {
				return 1
			}
		}
	}
	return 0
}

func requests_here(elev Elevator) int {
	for btn := 0; btn < driver.N_BUTTONS; btn++ {
		if intToBool(elev.Request[elev.Floor][btn]) {
			return 1
		}
	}
	return 0
}

func Requests_chooseDirection(elev Elevator) DirnBehaviourPair {
	switch elev.Dirn {
	case driver.MD_Up:
		if intToBool(requests_above(elev)) {
			return DirnBehaviourPair{driver.MD_Up, EB_Moving}
		}
		if intToBool(requests_here(elev)) {
			return DirnBehaviourPair{driver.MD_Down, EB_DoorOpen}
		}
		if intToBool(requests_below(elev)) {
			return DirnBehaviourPair{driver.MD_Down, EB_Moving}
		}
		return DirnBehaviourPair{driver.MD_Stop, EB_Idle}
	case driver.MD_Down:
		if intToBool(requests_below(elev)) {
			return DirnBehaviourPair{driver.MD_Down, EB_Moving}
		}
		if intToBool(requests_here(elev)) {
			return DirnBehaviourPair{driver.MD_Up, EB_DoorOpen}
		}
		if intToBool(requests_above(elev)) {
			return DirnBehaviourPair{driver.MD_Up, EB_Moving}
		}
		return DirnBehaviourPair{driver.MD_Stop, EB_Idle}
	case driver.MD_Stop:
		if intToBool(requests_here(elev)) {
			return DirnBehaviourPair{driver.MD_Stop, EB_DoorOpen}
		}
		if intToBool(requests_above(elev)) {
			return DirnBehaviourPair{driver.MD_Up, EB_Moving}
		}
		if intToBool(requests_below(elev)) {
			return DirnBehaviourPair{driver.MD_Down, EB_Moving}
		}
		return DirnBehaviourPair{driver.MD_Stop, EB_Idle}
	default:
		return DirnBehaviourPair{driver.MD_Stop, EB_Idle}
	}
}

func Requests_shouldStop(elev Elevator) bool {
	switch elev.Dirn {
	case driver.MD_Down:
		return (intToBool(elev.Request[elev.Floor][driver.BT_HallDown])) ||
			(intToBool(elev.Request[elev.Floor][driver.BT_Cab])) ||
			!intToBool(requests_below(elev))
	case driver.MD_Up:
		return (intToBool(elev.Request[elev.Floor][driver.BT_HallUp])) ||
			(intToBool(elev.Request[elev.Floor][driver.BT_Cab])) ||
			!intToBool(requests_above(elev))
	case driver.MD_Stop:
		return true
	default:
		return true
	}
}

func Requests_shouldClearImmediately(elev Elevator, btn_floor int, btn_type driver.ButtonType) bool {
	return (elev.Floor == btn_floor &&
		((elev.Dirn == driver.MD_Up && btn_type == driver.BT_HallUp) ||
			(elev.Dirn == driver.MD_Down && btn_type == driver.BT_HallDown) ||
			(elev.Dirn == driver.MD_Stop) ||
			(btn_type == driver.BT_Cab)))
}

func Requests_clearAtCurrentFloor(elev *Elevator, myIP string, upd_request chan<- UpdateRequest) {
	
	elev.SetElevatorRequest(elev.Floor, driver.BT_Cab, 0)
	upd_request<- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: elev.Floor, Button: driver.BT_Cab})

	switch elev.Dirn {
	case driver.MD_Up:

		if !intToBool(requests_above(*elev)) && !intToBool(elev.GetElevatorRequest(elev.Floor, int(driver.BT_HallUp))) {
			elev.SetElevatorRequest(elev.Floor, driver.BT_HallDown, 0)
			upd_request<- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: elev.Floor, Button: driver.BT_HallDown})
		}
		elev.SetElevatorRequest(elev.Floor, int(driver.BT_HallUp), 0)
		upd_request<- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: elev.Floor, Button: driver.BT_HallUp})


	case driver.MD_Down:

		if !intToBool(requests_below(*elev)) && !intToBool(elev.GetElevatorRequest(elev.Floor, int(driver.BT_HallDown))) {
			elev.SetElevatorRequest(elev.Floor, int(driver.BT_HallUp), 0)
			upd_request<- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: elev.Floor, Button: driver.BT_HallUp})
		}
		elev.SetElevatorRequest(elev.Floor, driver.BT_HallDown, 0)
		upd_request<- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: elev.Floor, Button: driver.BT_HallDown})


	case driver.MD_Stop:

		elev.SetElevatorRequest(elev.Floor, int(driver.BT_HallUp), 0)
		elev.SetElevatorRequest(elev.Floor, driver.BT_HallDown, 0)
		upd_request<- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: elev.Floor, Button: driver.BT_HallUp})
		upd_request<- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: elev.Floor, Button: driver.BT_HallDown})


	default:

		elev.SetElevatorRequest(elev.Floor, int(driver.BT_HallUp), 0)
		elev.SetElevatorRequest(elev.Floor, driver.BT_HallDown, 0)
		upd_request<- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: elev.Floor, Button: driver.BT_HallUp})
		upd_request<- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: elev.Floor, Button: driver.BT_HallDown})
	}
}

func intToBool(a int) bool {
	return a != 0
}
