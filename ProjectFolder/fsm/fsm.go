package fsm

import (
	"Sanntid/elevator"
	"Sanntid/elevator_io"
	"Sanntid/elevio"
)

var elev elevator.Elevator
var outputDevice elevator_io.ElevOutputDevice

func init() {
	elev = elevator.Elevator_uninitialized()
	outputDevice = elevator_io.Elevio_getOutputDevice()
}

func setAllLights(es elevator.Elevator) {
	for floor := 0; floor < elevator_io.N_FLOORS; floor++ {
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			outputDevice.RequestButtonLight(floor, elevio.ButtonType(btn), es.requests[floor][btn])
		}
	}
}