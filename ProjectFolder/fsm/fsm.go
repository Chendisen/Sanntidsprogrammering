package fsm

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/elevator_io"
	"fmt"
	"runtime"
	"Sanntid/requests"
	"Sanntid/timer"
)

func toBool(a int) bool {
	return a != 0
}

var elev elevator.Elevator
var outputDevice elevator_io.ElevOutputDevice

func init() {
	elev = elevator.Elevator_uninitialized()
	outputDevice = elevator_io.Elevio_getOutputDevice()
}

func setAllLights(es elevator.Elevator) {

	for floor := 0; floor < elevator_io.N_FLOORS; floor++ {
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			outputDevice.RequestButtonLight(floor, driver.ButtonType(btn), toBool(es.Request[floor][btn]))
		}
	}
}

func fsm_onInitBetweenFloors() {
	outputDevice.MotorDirection(driver.MD_Down)
	elev.Dirn = driver.MD_Down
	elev.Behaviour = elevator.EB_Moving
}

func fsm_onRequestButtonPress(btn_floor int, btn_type driver.ButtonType) {

	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, elevator_io.Elevio_button_toString(btn_type))

	elevator.Elevator_print(elev)

	switch(elev.Behaviour){

	case elevator.EB_DoorOpen:
		if(requests.Requests_shouldClearImmediately(elev, btn_floor, btn_type)) {
			timer.Timer_start(elev.Config.DoorOpenDuration_s)
		} else {
			elev.Request[btn_floor][int(btn_type)] = 1
		}
		break
	
	case elevator.EB_Moving:
		elev.Request[btn_floor][int(btn_type)] = 1
		break

	case elevator.EB_Idle:
		
		elev.Request[btn_floor][int(btn_type)] = 1
		pair := requests.Requests_chooseDirection(elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour


	}
}
