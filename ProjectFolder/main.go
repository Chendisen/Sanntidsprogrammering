package main

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/fsm"
	"Sanntid/timer"
	"Sanntid/world_view"
	"Sanntid/network/localip"
	"fmt"
)

func main() {

	const numFloors int = 4
	driver.Init("localhost:15657", numFloors)

	var elev elevator.Elevator = elevator.Elevator_uninitialized()
	var tmr timer.Timer = timer.Timer_uninitialized()
	var alv_list world_view.AliveList = world_view.MakeAliveList()
	var wld_view world_view.WorldView = world_view.MakeWorldView(alv_list.MyIP)	

	//var d driver.MotorDirection = driver.MD_Up
	//driver.SetMotorDirection(d)

	drv_buttons := make(chan driver.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	ord_updated := make(chan int)

	go driver.PollButtons(drv_buttons)
	go driver.PollFloorSensor(drv_floors)
	go driver.PollObstructionSwitch(drv_obstr)
	go driver.PollStopButton(drv_stop)
	go fsm.Fsm_checkTimeOut(&elev, &wld_view, alv_list.myIP, &tmr)

	// a:= <- drv_floors
	// if a==-1 {
	// 	fsm.Fsm_onInitBetweenFloors(&elevator)
	// }

	fsm.Fsm_onInitBetweenFloors(&elev)

	for {
		select {
		case a := <-drv_buttons:

			select {
			case a.Button == 2:
				Fsm_onRequestButtonPress(&elev, &wld_view, alv_list.MyIP, &tmr, a.Floor, a.Button)
			default: 
				wld_view.SetHallRequestAtFloor(a.Floor, a.Button)
			}

			// Press of button shall update my worldview which will then propagate out and be published that new info has been found.
			// 		But we must seperate between cab and hall buttons since cab calls can only be handled by itself. 
			// We must then have an own function for reading in the world view and update the requests matrix of the elevator. 
			// Must find out how to make elevator only considers its requests matrix and state to make decisions.  

			// When we change state we also need to update the world view, this should probably happen when we arrive at new floor and update state.


			fmt.Printf("%+v\n", a)
			//driver.SetButtonLamp(a.Button, a.Floor, true)
			// fsm.Fsm_onRequestButtonPress(&elev, &tmr, a.Floor, a.Button)
			fmt.Printf("Request floor: %d", a.Floor)

		case a := <-drv_floors:
			fmt.Printf("This floor polled: %d\n", a)
			// if a == numFloors-1 {
			//     d = driver.MD_Down
			// } else if a == 0 {
			//     d = driver.MD_Up
			// }
			// driver.SetMotorDirection(d)
			fsm.Fsm_onFloorArrival(&elev, &wld_view, alv_list.MyIP, &tmr, a)

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
		case a := <-ord_updated:
			for floor, buttons := range wld_view.GetMyAssignedOrders() {
				for button, value := range buttons {
					if value == True {
						Fsm_onRequestButtonPress(&elev, &wld_view, alv_list.MyIP, &tmr, a.Floor, a.Button)
					}
					else{
						es.Request[floor][button] = 0
					}
				}
			}
		}
	}
}

// TODO: Floor sensor lights dont work properly
// TODO: Elevator picks up people going in both directions when entering floor

// Help me