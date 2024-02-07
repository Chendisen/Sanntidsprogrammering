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

var elev elevator.Elevator
var outputDevice elevator_io.ElevOutputDevice

func init() {
	elev = elevator.Elevator_uninitialized()
	outputDevice = elevator_io.Elevio_getOutputDevice()
}

func setAllLights(es elevator.Elevator) {

	for floor := 0; floor < elevator_io.N_FLOORS; floor++ {
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			outputDevice.RequestButtonLight(floor, driver.ButtonType(btn), driver.IntToBool(es.Request[floor][btn]))
		}
	}
}

func Fsm_onInitBetweenFloors() {
	outputDevice.MotorDirection(driver.MD_Down)
	elev.Dirn = driver.MD_Down
	elev.Behaviour = elevator.EB_Moving
}

func Fsm_onRequestButtonPress(btn_floor int, btn_type driver.ButtonType) {

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
			
	case elevator.EB_Moving:
		elev.Request[btn_floor][int(btn_type)] = 1
		
	case elevator.EB_Idle:

		elev.Request[btn_floor][int(btn_type)] = 1
		pair := requests.Requests_chooseDirection(elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour

		switch(pair.Behaviour){

		case elevator.EB_DoorOpen:
			outputDevice.DoorLight(true)
			timer.Timer_start(elev.Config.DoorOpenDuration_s)
			elev = requests.Requests_clearAtCurrentFloor(elev)
			
		case elevator.EB_Moving:
			outputDevice.MotorDirection(elev.Dirn)
			
		case elevator.EB_Idle:
					}

			}

	setAllLights(elev)

	fmt.Printf("\nNew state:\n")
	elevator.Elevator_print(elev)
}

func Fsm_onFloorArrival(newFloor int) {
	pc, _, _, _ := runtime.Caller(0) 
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d)\n", functionName, newFloor) //uuuuuhhhm what is all this
	elevator.Elevator_print(elev)

	elev.Floor = newFloor

	outputDevice.FloorIndicator(elev.Floor)

	switch(elev.Behaviour){
	case elevator.EB_Moving:
		if(requests.Requests_shouldStop(elev)){
			outputDevice.MotorDirection(driver.MD_Stop)
			outputDevice.DoorLight(true)
			elev= requests.Requests_clearAtCurrentFloor(elev)
			timer.Timer_start(elev.Config.DoorOpenDuration_s)
			setAllLights(elev)
			elev.Behaviour = elevator.EB_DoorOpen
		}
	default:
	}

	fmt.Printf("\nNew State:\n")
	elevator.Elevator_print(elev)
}

func Fsm_onDoorTimeout() {
	pc, _, _, _ := runtime.Caller(0) 
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s()\n", functionName) //uuuuuhhhm what is all this
	elevator.Elevator_print(elev)

	switch(elev.Behaviour){
	case elevator.EB_DoorOpen:
		pair := requests.Requests_chooseDirection(elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour

		switch(elev.Behaviour){
		case elevator.EB_DoorOpen:
			timer.Timer_start(elev.Config.DoorOpenDuration_s)
			elev = requests.Requests_clearAtCurrentFloor(elev)
			setAllLights(elev)

		case elevator.EB_Moving:
			outputDevice.DoorLight(false)
			outputDevice.MotorDirection(elev.Dirn)
		case elevator.EB_Idle:
			outputDevice.DoorLight(false)
			outputDevice.MotorDirection(elev.Dirn)
		}
	default:
	}

	fmt.Printf("\nNew State:\n")
	elevator.Elevator_print(elev)

}