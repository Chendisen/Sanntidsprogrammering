package fsm

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/timer"
	"fmt"
	"runtime"
)

func Fsm_onInitBetweenFloors(elev *elevator.Elevator, myIP string, set_behaviour chan<- elevator.ElevatorBehaviour , set_direction chan<- driver.MotorDirection) {
	if elev.Floor == 0 {
		driver.SetMotorDirection(driver.MD_Up)
		elev.Dirn = driver.MD_Up
		set_direction<- driver.MD_Up
	} else {
		driver.SetMotorDirection(driver.MD_Down)
		elev.Dirn = driver.MD_Down
		set_direction<- driver.MD_Down
	}

	elev.Behaviour = elevator.EB_Moving
	set_behaviour<- elevator.EB_Moving
}

func Fsm_onRequestButtonPress(elev *elevator.Elevator, myIP string, tmr *timer.Timer, watchdog *timer.Timer, btn_floor int, btn_type driver.ButtonType, set_behaviour chan<- elevator.ElevatorBehaviour, set_direction chan<- driver.MotorDirection, finished_request_at_floor chan<- driver.ButtonEvent) {

	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, driver.Driver_button_toString(btn_type))

	switch elev.Behaviour {

	case elevator.EB_DoorOpen:
		if Requests_shouldClearImmediately(*elev, btn_floor, btn_type) {
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
			finished_request_at_floor<- driver.ButtonEvent{Floor: btn_floor, Button: btn_type}
		} else {
			elev.SetElevatorRequest(btn_floor, int(btn_type), 1)
		}

	case elevator.EB_Moving:
		elev.SetElevatorRequest(btn_floor, int(btn_type), 1)

	case elevator.EB_Idle:

		elev.SetElevatorRequest(btn_floor, int(btn_type), 1)
		pair := Requests_chooseDirection(*elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour
		set_direction<- pair.Dirn
		set_behaviour<- pair.Behaviour

		switch pair.Behaviour {
		case elevator.EB_DoorOpen:
			driver.SetDoorOpenLamp(true)
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
			Requests_clearAtCurrentFloor(elev, myIP, finished_request_at_floor)

		case elevator.EB_Moving:
			driver.SetMotorDirection(elev.Dirn)
			watchdog.Timer_start(timer.WATCHDOG_TimeoutTime)

		case elevator.EB_Idle:
		}
	}
}

func Fsm_onFloorArrival(elev *elevator.Elevator, myIP string, tmr *timer.Timer, newFloor int, set_behaviour chan<- elevator.ElevatorBehaviour, set_floor chan<- int, finished_request_at_floor chan<- driver.ButtonEvent) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d)\n", functionName, newFloor) //uuuuuhhhm what is all this

	elev.Floor = newFloor
	set_floor<- newFloor

	driver.SetFloorIndicator(elev.Floor)

	switch elev.Behaviour {
	case elevator.EB_Moving:
		if Requests_shouldStop(*elev) {

			driver.SetMotorDirection(driver.MD_Stop)
			driver.SetDoorOpenLamp(true)

			Requests_clearAtCurrentFloor(elev, myIP, finished_request_at_floor)
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)

			elev.Behaviour = elevator.EB_DoorOpen
			set_behaviour<- elevator.EB_DoorOpen
		}
	default:
	}
}

func Fsm_onDoorTimeout(elev *elevator.Elevator, myIP string, tmr *timer.Timer, watchdog *timer.Timer, set_behaviour chan<- elevator.ElevatorBehaviour, set_direction chan<- driver.MotorDirection, finished_request_at_floor chan<- driver.ButtonEvent) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s()\n", functionName) //uuuuuhhhm what is all this

	switch elev.Behaviour {
	case elevator.EB_DoorOpen:
		pair := Requests_chooseDirection(*elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour
		set_direction<- pair.Dirn
		set_behaviour<- pair.Behaviour

		switch elev.Behaviour {
		case elevator.EB_DoorOpen:
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
			Requests_clearAtCurrentFloor(elev, myIP, finished_request_at_floor)

		case elevator.EB_Moving:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(elev.Dirn)
			watchdog.Timer_start(timer.WATCHDOG_TimeoutTime)

		case elevator.EB_Idle:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(elev.Dirn)
			watchdog.Timer_start(timer.WATCHDOG_TimeoutTime)
		}
	default:
	}
}
