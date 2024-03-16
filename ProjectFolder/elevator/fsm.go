package elevator

import (
	"Sanntid/resources/driver"
	"Sanntid/timer"
	. "Sanntid/resources/update_request"
	"fmt"
	"runtime"
)

func Fsm_onInitBetweenFloors(elev *Elevator, myIP string, upd_request chan UpdateRequest) {
	if elev.Floor == 0 {
		driver.SetMotorDirection(driver.MD_Up)
		elev.Dirn = driver.MD_Up
		upd_request <- GenerateUpdateRequest(SetDirection, driver.MD_Up)
	} else {
		driver.SetMotorDirection(driver.MD_Down)
		elev.Dirn = driver.MD_Down
		upd_request <- GenerateUpdateRequest(SetDirection, driver.MD_Down)
	}

	elev.Behaviour = EB_Moving
	upd_request <- GenerateUpdateRequest(SetBehaviour, EB_Moving)
}

func Fsm_onRequestButtonPress(elev *Elevator, myIP string, tmr *timer.Timer, watchdog *timer.Timer, btn_floor int, btn_type driver.ButtonType, upd_request chan UpdateRequest) {

	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, driver.DriverButtonToString(btn_type))

	switch elev.Behaviour {

	case EB_DoorOpen:
		fmt.Println("Door open")
		if Requests_shouldClearImmediately(*elev, btn_floor, btn_type) {
			fmt.Println("Req should clear immediately")
			tmr.TimerStart(elev.Config.DoorOpenDuration_s)
			upd_request <- GenerateUpdateRequest(FinishedRequestAtFloor, driver.ButtonEvent{Floor: btn_floor, Button: btn_type})
		} else {
			elev.SetElevatorRequest(btn_floor, int(btn_type), 1)
		}

	case EB_Moving:
		
		elev.SetElevatorRequest(btn_floor, int(btn_type), 1)

	case EB_Idle:

		elev.SetElevatorRequest(btn_floor, int(btn_type), 1)
		pair := Requests_chooseDirection(*elev)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour
		upd_request <- GenerateUpdateRequest(SetDirection, pair.Dirn)
		upd_request <- GenerateUpdateRequest(SetBehaviour, pair.Behaviour)

		switch pair.Behaviour {
		case EB_DoorOpen:
			driver.SetDoorOpenLamp(true)
			tmr.TimerStart(elev.Config.DoorOpenDuration_s)
			Requests_clearAtCurrentFloor(elev, myIP, upd_request)

		case EB_Moving:
			driver.SetMotorDirection(elev.Dirn)
			watchdog.TimerStart(timer.WATCHDOG_TimeoutTime)

		case EB_Idle:
		}
	}
}

func Fsm_initAllOrders(ord_updated chan<- bool) {
	ord_updated<- true
}

func Fsm_setAssignedOrders(assignedOrders [][2]bool, elev *Elevator, myIP string, timerDoor *timer.Timer, timerWatchdog *timer.Timer, upd_request chan UpdateRequest) {
	for floor, buttons := range assignedOrders {
		for button, value := range buttons {
			if value {
				Fsm_onRequestButtonPress(elev, myIP, timerDoor, timerWatchdog, floor, driver.ButtonType(button), upd_request)

			} else {
				elev.ClearElevatorLight(floor, button)
			}
		}
	}
}

func Fsm_setCabOrders(cabRequests []bool, elev *Elevator, myIP string, timerDoor *timer.Timer, timerWatchdog *timer.Timer, upd_request chan UpdateRequest) {
	for floor, value := range cabRequests {
		if value {
			Fsm_onRequestButtonPress(elev, myIP, timerDoor, timerWatchdog, floor, driver.BT_Cab, upd_request)
		} else {
			elev.ClearElevatorLight(floor, driver.BT_Cab)
		}
	}
}

func Fsm_onFloorArrival(elev *Elevator, myIP string, tmr *timer.Timer, newFloor int, upd_request chan UpdateRequest) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d)\n", functionName, newFloor) //uuuuuhhhm what is all this

	elev.Floor = newFloor
	upd_request <- GenerateUpdateRequest(SetFloor, newFloor)
	driver.SetFloorIndicator(elev.Floor)

	switch elev.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(*elev) {

			driver.SetMotorDirection(driver.MD_Stop)
			driver.SetDoorOpenLamp(true)

			Requests_clearAtCurrentFloor(elev, myIP, upd_request)
			tmr.TimerStart(elev.Config.DoorOpenDuration_s)
			fmt.Println("Elevator should stop here")

			elev.Behaviour = EB_DoorOpen
			upd_request <- GenerateUpdateRequest(SetBehaviour, EB_DoorOpen)
		}
	default:
	}
}

func Fsm_onDoorTimeout(elev *Elevator, myIP string, tmr *timer.Timer, watchdog *timer.Timer, upd_request chan UpdateRequest) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s()\n", functionName) //uuuuuhhhm what is all this

	switch elev.Behaviour {
	case EB_DoorOpen:
		pair := Requests_chooseDirection(*elev)
		fmt.Printf("The old direction of the elevator is: %d", elev.Dirn)
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour

		upd_request <- GenerateUpdateRequest(SetDirection, elev.Dirn)
		upd_request <- GenerateUpdateRequest(SetBehaviour, elev.Behaviour)

		switch elev.Behaviour {
		case EB_DoorOpen:
			fmt.Println("The door is still open")
			tmr.TimerStart(elev.Config.DoorOpenDuration_s)
			Requests_clearAtCurrentFloor(elev, myIP, upd_request)

		case EB_Moving:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(elev.Dirn)
			watchdog.TimerStart(timer.WATCHDOG_TimeoutTime)

		case EB_Idle:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(elev.Dirn)
			watchdog.TimerStart(timer.WATCHDOG_TimeoutTime)
		}
	default:
	}
}
