package fsm

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/requests"
	"Sanntid/timer"
	"Sanntid/world_view"
	"fmt"
	"runtime"
)

// var elev elevator.Elevator
// var outputDevice elevator_io.ElevOutputDevice

// func Init() {
// 	elev = elevator.Elevator_uninitialized()
// 	outputDevice = elevator_io.Elevio_getOutputDevice()
// }

func SetAllLights(lightArray [][3]bool) {
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			//outputDevice.RequestButtonLight(floor, driver.ButtonType(btn), driver.IntToBool(es.Request[floor][btn]))
			driver.SetButtonLamp(driver.ButtonType(btn), floor, lightArray[floor][btn])
		}
	}
}

func Fsm_onInitBetweenFloors(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string) {
	driver.SetMotorDirection(driver.MD_Down)
	es.Dirn = driver.MD_Down
	wld_view.SetDirection(myIP, driver.MD_Down)
	es.Behaviour = elevator.EB_Moving
	wld_view.SetBehaviour(myIP, elevator.EB_Moving)
}

func Fsm_onRequestButtonPress(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string, tmr *timer.Timer, btn_floor int, btn_type driver.ButtonType) {

	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, driver.Driver_button_toString(btn_type))

	//elevator.Elevator_print(*es)

	switch es.Behaviour {

	case elevator.EB_DoorOpen:
		if requests.Requests_shouldClearImmediately(*es, btn_floor, btn_type) {
			timer.Timer_start(tmr, es.Config.DoorOpenDuration_s)
			wld_view.FinishedRequestAtFloor(myIP, btn_floor, btn_type)
		} else {
			es.Request[btn_floor][int(btn_type)] = 1
		}

	case elevator.EB_Moving:
		es.Request[btn_floor][int(btn_type)] = 1

	case elevator.EB_Idle:

		es.Request[btn_floor][int(btn_type)] = 1
		pair := requests.Requests_chooseDirection(*es)
		es.Dirn = pair.Dirn
		wld_view.SetDirection(myIP, pair.Dirn)
		es.Behaviour = pair.Behaviour
		wld_view.SetBehaviour(myIP, pair.Behaviour)

		switch pair.Behaviour {

		case elevator.EB_DoorOpen:
			driver.SetDoorOpenLamp(true)
			timer.Timer_start(tmr, es.Config.DoorOpenDuration_s)
			requests.Requests_clearAtCurrentFloor(es, wld_view, myIP)

		case elevator.EB_Moving:
			driver.SetMotorDirection(es.Dirn)

		case elevator.EB_Idle:
		}

	}

	//fmt.Printf("\nNew state:\n")
	//elevator.Elevator_print(*es)
}

func Fsm_onFloorArrival(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string, tmr *timer.Timer, newFloor int) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d)\n", functionName, newFloor) //uuuuuhhhm what is all this
	//elevator.Elevator_print(*es)

	es.Floor = newFloor
	wld_view.SetFloor(myIP, newFloor)

	driver.SetFloorIndicator(es.Floor)

	switch es.Behaviour {
	case elevator.EB_Moving:
		if requests.Requests_shouldStop(*es) {

			es.Dirn = driver.MD_Stop
			wld_view.SetDirection(myIP, driver.MD_Stop)
			driver.SetMotorDirection(es.Dirn)
			driver.SetDoorOpenLamp(true)

			requests.Requests_clearAtCurrentFloor(es, wld_view, myIP)

			timer.Timer_start(tmr, es.Config.DoorOpenDuration_s)

			es.Behaviour = elevator.EB_DoorOpen
			wld_view.SetBehaviour(myIP, elevator.EB_DoorOpen)
		}
	default:
	}

	//fmt.Printf("\nNew State:\n")
	//elevator.Elevator_print(*es)
}

func Fsm_onDoorTimeout(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string, tmr *timer.Timer) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s()\n", functionName) //uuuuuhhhm what is all this
	//elevator.Elevator_print(*es)

	switch es.Behaviour {
	case elevator.EB_DoorOpen:
		pair := requests.Requests_chooseDirection(*es)
		es.Dirn = pair.Dirn
		wld_view.SetDirection(myIP, pair.Dirn)
		es.Behaviour = pair.Behaviour
		wld_view.SetBehaviour(myIP, pair.Behaviour)

		switch es.Behaviour {
		case elevator.EB_DoorOpen:
			timer.Timer_start(tmr, es.Config.DoorOpenDuration_s)
			requests.Requests_clearAtCurrentFloor(es, wld_view, myIP)

		case elevator.EB_Moving:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(es.Dirn)

		case elevator.EB_Idle:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(es.Dirn)
		}
	default:
	}

	//fmt.Printf("\nNew State:\n")
	//elevator.Elevator_print(*es)

}

func Fsm_checkTimeOut(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string, tmr *timer.Timer) {
	for {
		if es.DoorObstructed {
			timer.Timer_start(tmr, es.Config.DoorOpenDuration_s)
		}
		if timer.Timer_timedOut(tmr) {
			timer.Timer_stop(tmr)
			Fsm_onDoorTimeout(es, wld_view, myIP, tmr)
		}
	}
}
