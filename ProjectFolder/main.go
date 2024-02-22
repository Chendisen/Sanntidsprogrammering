package main

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/fsm"
	"Sanntid/timer"
	"fmt"
)

func main() {

	numFloors := 4
	driver.Init("localhost:15657", numFloors)

	var elev elevator.Elevator = elevator.Elevator_uninitialized()
	var tmr timer.Timer = timer.Timer_uninitialized()

	//var d driver.MotorDirection = driver.MD_Up
	//driver.SetMotorDirection(d)

	drv_buttons := make(chan driver.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go driver.PollButtons(drv_buttons)
	go driver.PollFloorSensor(drv_floors)
	go driver.PollObstructionSwitch(drv_obstr)
	go driver.PollStopButton(drv_stop)
	go fsm.Fsm_checkTimeOut(&elev, &tmr)

	// a:= <- drv_floors
	// if a==-1 {
	// 	fsm.Fsm_onInitBetweenFloors(&elevator)
	// }

	fsm.Fsm_onInitBetweenFloors(&elev)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			//driver.SetButtonLamp(a.Button, a.Floor, true)
			fsm.Fsm_onRequestButtonPress(&elev, &tmr, a.Floor, a.Button)
			fmt.Printf("Request floor: %d", a.Floor)

		case a := <-drv_floors:
			fmt.Printf("This floor polled: %d\n", a)
			// if a == numFloors-1 {
			//     d = driver.MD_Down
			// } else if a == 0 {
			//     d = driver.MD_Up
			// }
			// driver.SetMotorDirection(d)
			fsm.Fsm_onFloorArrival(&elev, &tmr, a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a && elev.Behaviour == elevator.EB_DoorOpen {
				elev.DoorObstructed = true
			} else {
				elev.DoorObstructed = false
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := driver.ButtonType(0); b < 3; b++ {
					driver.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}

// TODO: Floor sensor lights dont work properly
// TODO: Elevator picks up people going in both directions when entering floor
