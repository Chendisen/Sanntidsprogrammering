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
		// set_direction<- driver.MD_Up
		fmt.Println("We came from init1")
		upd_request <- GenerateUpdateRequest(SetDirection, driver.MD_Up)
	} else {
		driver.SetMotorDirection(driver.MD_Down)
		elev.Dirn = driver.MD_Down
		// set_direction<- driver.MD_Down
		fmt.Println("We came from init2")
		upd_request <- GenerateUpdateRequest(SetDirection, driver.MD_Down)
	}

	elev.Behaviour = EB_Moving
	// set_behaviour<- EB_Moving
	upd_request <- GenerateUpdateRequest(SetBehaviour, EB_Moving)
}

func Fsm_onRequestButtonPress(elev *Elevator, myIP string, tmr *timer.Timer, watchdog *timer.Timer, btn_floor int, btn_type driver.ButtonType, upd_request chan UpdateRequest) {

	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d, %s)\n", functionName, btn_floor, driver.DriverButtonToString(btn_type))

	switch elev.Behaviour {

	case EB_DoorOpen:
		if Requests_shouldClearImmediately(*elev, btn_floor, btn_type) {
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
			//finished_request_at_floor<- driver.ButtonEvent{Floor: btn_floor, Button: btn_type}
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
		//set_direction<- pair.Dirn
		//set_behaviour<- pair.Behaviour
		fmt.Println("We came from req1")
		upd_request <- GenerateUpdateRequest(SetDirection, pair.Dirn)
		upd_request <- GenerateUpdateRequest(SetBehaviour, pair.Behaviour)

		switch pair.Behaviour {
		case EB_DoorOpen:
			driver.SetDoorOpenLamp(true)
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
			Requests_clearAtCurrentFloor(elev, myIP, upd_request)

		case EB_Moving:
			driver.SetMotorDirection(elev.Dirn)
			watchdog.Timer_start(timer.WATCHDOG_TimeoutTime)

		case EB_Idle:
		}
	}
}

func Fsm_onFloorArrival(elev *Elevator, myIP string, tmr *timer.Timer, newFloor int, upd_request chan UpdateRequest) {
	pc, _, _, _ := runtime.Caller(0)
	functionName := runtime.FuncForPC(pc).Name()

	fmt.Printf("\n\n%s(%d)\n", functionName, newFloor) //uuuuuhhhm what is all this

	elev.Floor = newFloor
	//set_floor<- newFloor
	upd_request <- GenerateUpdateRequest(SetFloor, newFloor)
	driver.SetFloorIndicator(elev.Floor)

	switch elev.Behaviour {
	case EB_Moving:
		if Requests_shouldStop(*elev) {

			driver.SetMotorDirection(driver.MD_Stop)
			driver.SetDoorOpenLamp(true)

			Requests_clearAtCurrentFloor(elev, myIP, upd_request)
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)

			elev.Behaviour = EB_DoorOpen
			//set_behaviour<- EB_DoorOpen
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
		elev.Dirn = pair.Dirn
		elev.Behaviour = pair.Behaviour
		//set_direction<- pair.Dirn
		//set_behaviour<- pair.Behaviour

		fmt.Printf("The type of the data is: %T\n", elev.Dirn)
		fmt.Println("We came from req1")
		upd_request <- GenerateUpdateRequest(SetDirection, elev.Dirn)

		fmt.Printf("We sent the direction")
		upd_request <- GenerateUpdateRequest(SetBehaviour, elev.Behaviour)

		switch elev.Behaviour {
		case EB_DoorOpen:
			tmr.Timer_start(elev.Config.DoorOpenDuration_s)
			Requests_clearAtCurrentFloor(elev, myIP, upd_request)

		case EB_Moving:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(elev.Dirn)
			watchdog.Timer_start(timer.WATCHDOG_TimeoutTime)

		case EB_Idle:
			driver.SetDoorOpenLamp(false)
			driver.SetMotorDirection(elev.Dirn)
			watchdog.Timer_start(timer.WATCHDOG_TimeoutTime)
		}
	default:
	}
}
