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

const watchdogTime float64 = 3

func Fsm_onInitBetweenFloors(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string) {
	if(es.Floor == 0){
		driver.SetMotorDirection(driver.MD_Up)
		es.Dirn = driver.MD_Up
		wld_view.SetDirection(myIP, driver.MD_Up)
	} else {
		driver.SetMotorDirection(driver.MD_Down)
		es.Dirn = driver.MD_Down
		wld_view.SetDirection(myIP, driver.MD_Down)
	}
	
	es.Behaviour = elevator.EB_Moving
	wld_view.SetBehaviour(myIP, elevator.EB_Moving)
}

func Fsm_onRequestButtonPress(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string, tmr *timer.Timer, watchdog *timer.Timer, btn_floor int, btn_type driver.ButtonType) {

	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, driver.Driver_button_toString(btn_type))

	switch es.Behaviour {

	case elevator.EB_DoorOpen:
		if requests.Requests_shouldClearImmediately(*es, btn_floor, btn_type) {
			tmr.Timer_start(es.Config.DoorOpenDuration_s)
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
			tmr.Timer_start(es.Config.DoorOpenDuration_s)
			requests.Requests_clearAtCurrentFloor(es, wld_view, myIP)

		case elevator.EB_Moving:
			driver.SetMotorDirection(es.Dirn)
			watchdog.Timer_start(watchdogTime)

		case elevator.EB_Idle:
		}
	}
}

func Fsm_onFloorArrival(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string, tmr *timer.Timer, newFloor int) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d)\n", functionName, newFloor) //uuuuuhhhm what is all this

	es.Floor = newFloor
	wld_view.SetFloor(myIP, newFloor)

	driver.SetFloorIndicator(es.Floor)

	switch es.Behaviour {
	case elevator.EB_Moving:
		if requests.Requests_shouldStop(*es) {

			driver.SetMotorDirection(driver.MD_Stop)
			driver.SetDoorOpenLamp(true)

			requests.Requests_clearAtCurrentFloor(es, wld_view, myIP)

			tmr.Timer_start(es.Config.DoorOpenDuration_s)

			es.Behaviour = elevator.EB_DoorOpen
			wld_view.SetBehaviour(myIP, elevator.EB_DoorOpen)
		}
	default:
	}
}

func Fsm_onDoorTimeout(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string, tmr *timer.Timer, watchdog *timer.Timer) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s()\n", functionName) //uuuuuhhhm what is all this

	switch es.Behaviour {
	case elevator.EB_DoorOpen:
		pair := requests.Requests_chooseDirection(*es)
		es.Dirn = pair.Dirn
		wld_view.SetDirection(myIP, pair.Dirn)
		es.Behaviour = pair.Behaviour
		wld_view.SetBehaviour(myIP, pair.Behaviour)

		switch es.Behaviour {
		case elevator.EB_DoorOpen:
			tmr.Timer_start(es.Config.DoorOpenDuration_s)
			requests.Requests_clearAtCurrentFloor(es, wld_view, myIP)

		case elevator.EB_Moving:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(es.Dirn)
			watchdog.Timer_start(watchdogTime)

		case elevator.EB_Idle:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(es.Dirn)
			watchdog.Timer_start(watchdogTime)
		}
	default:
	}
}

func Fsm_checkTimeOut(es *elevator.Elevator, wld_view *world_view.WorldView, myIP string, tmr *timer.Timer, watchdog *timer.Timer) {
	for {
		if es.DoorObstructed {
			tmr.Timer_start(es.Config.DoorOpenDuration_s)
		}
		if tmr.Timer_timedOut(es.Config.DoorOpenDuration_s) {
			tmr.Timer_stop()
			Fsm_onDoorTimeout(es, wld_view, myIP, tmr, watchdog)
		}
	}
}
