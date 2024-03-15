package fsm

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	. "Sanntid/resources"
)

type DirnBehaviourPair struct {
	Dirn      driver.MotorDirection
	Behaviour elevator.ElevatorBehaviour
}

func requests_above(e elevator.Elevator) int {
	for floor := e.Floor + 1; floor < driver.N_FLOORS; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			if intToBool(e.Request[floor][btn]) {
				return 1
			}
		}
	}
	return 0
}

func requests_below(e elevator.Elevator) int {
	for floor := 0; floor < e.Floor; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			if intToBool(e.Request[floor][btn]) {
				return 1
			}
		}
	}
	return 0
}

func requests_here(e elevator.Elevator) int {
	for btn := 0; btn < driver.N_BUTTONS; btn++ {
		if intToBool(e.Request[e.Floor][btn]) {
			return 1
		}
	}
	return 0
}

func Requests_chooseDirection(e elevator.Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case driver.MD_Up:
		if intToBool(requests_above(e)) {
			return DirnBehaviourPair{driver.MD_Up, elevator.EB_Moving}
		}
		if intToBool(requests_here(e)) {
			return DirnBehaviourPair{driver.MD_Down, elevator.EB_DoorOpen}
		}
		if intToBool(requests_below(e)) {
			return DirnBehaviourPair{driver.MD_Down, elevator.EB_Moving}
		}
		return DirnBehaviourPair{driver.MD_Stop, elevator.EB_Idle}
	case driver.MD_Down:
		if intToBool(requests_below(e)) {
			return DirnBehaviourPair{driver.MD_Down, elevator.EB_Moving}
		}
		if intToBool(requests_here(e)) {
			return DirnBehaviourPair{driver.MD_Up, elevator.EB_DoorOpen}
		}
		if intToBool(requests_above(e)) {
			return DirnBehaviourPair{driver.MD_Up, elevator.EB_Moving}
		}
		return DirnBehaviourPair{driver.MD_Stop, elevator.EB_Idle}
	case driver.MD_Stop:
		if intToBool(requests_here(e)) {
			return DirnBehaviourPair{driver.MD_Stop, elevator.EB_DoorOpen}
		}
		if intToBool(requests_above(e)) {
			return DirnBehaviourPair{driver.MD_Up, elevator.EB_Moving}
		}
		if intToBool(requests_below(e)) {
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
		return (intToBool(e.Request[e.Floor][driver.BT_HallDown])) ||
			(intToBool(e.Request[e.Floor][driver.BT_Cab])) ||
			!intToBool(requests_below(e))
	case driver.MD_Up:
		return (intToBool(e.Request[e.Floor][driver.BT_HallUp])) ||
			(intToBool(e.Request[e.Floor][driver.BT_Cab])) ||
			!intToBool(requests_above(e))
	case driver.MD_Stop:
		return true
	default:
		return true
	}
}

func Requests_shouldClearImmediately(elev elevator.Elevator, btn_floor int, btn_type driver.ButtonType) bool {
	return (elev.Floor == btn_floor &&
		((elev.Dirn == driver.MD_Up && btn_type == driver.BT_HallUp) ||
			(elev.Dirn == driver.MD_Down && btn_type == driver.BT_HallDown) ||
			(elev.Dirn == driver.MD_Stop) ||
			(btn_type == driver.BT_Cab)))
}

func Requests_clearAtCurrentFloor(elev *elevator.Elevator, myIP string, upd_request chan<- UpdateRequest) {
	
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

		if !intToBool(requests_below(*elev)) && !intToBool(elev.GetElevatorRequest(elev.Floor, int(driver.BT_HallUp))) {
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
