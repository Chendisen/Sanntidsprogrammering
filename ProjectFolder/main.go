package main

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/fsm"
	"Sanntid/network"
	"Sanntid/order_assigner"
	"Sanntid/process_pair"
	"Sanntid/timer"
	"Sanntid/world_view"
	"fmt"
	"os"
	"os/exec"
)

func main() {
 
	const numFloors int = 4

	var elev elevator.Elevator = elevator.Elevator_uninitialized()
	var tmr timer.Timer = timer.Timer_uninitialized()
	var alv_list world_view.AliveList = world_view.MakeAliveList()
	var wld_view world_view.WorldView = world_view.MakeWorldView(alv_list.MyIP)
	var hrd_list world_view.HeardFromList = world_view.MakeHeardFromList(alv_list.MyIP)
	
	lgt_array := make([][3]bool, numFloors)

	drv_buttons := make(chan driver.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	ord_updated := make(chan bool, 10)
	wld_updated := make(chan bool, 10)

	startNew := make(chan bool)



	go process_pair.ProcessPair(alv_list.MyIP, &wld_view, &tmr, startNew)

	for range startNew{
		break
	}

	//var path string = "~/Documents/EddChris/Sanntidsprogrammering/ProjectFolder"
	path2,_ := os.Getwd()
	cmd := exec.Command("gnome-terminal", "--window", "--", "sh", "-c", "cd "+path2+" && go run main.go")

	fmt.Printf("Path: %s", path2)

	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start myself")
		panic(err)
	}

	driver.Init("localhost:15657", numFloors)


	go driver.PollButtons(drv_buttons)
	go driver.PollFloorSensor(drv_floors)
	go driver.PollObstructionSwitch(drv_obstr)
	go driver.PollStopButton(drv_stop)
	go fsm.Fsm_checkTimeOut(&elev, &wld_view, alv_list.MyIP, &tmr)
	go network.StartCommunication(alv_list.MyIP, &wld_view, &alv_list, &hrd_list, &lgt_array, ord_updated, wld_updated)

	fsm.Fsm_onInitBetweenFloors(&elev, &wld_view, alv_list.MyIP)
	world_view.InitLights(&lgt_array, alv_list.MyIP, wld_view)
	ord_updated<-true

	for {
		select {
		case a := <-drv_buttons:

			wld_view.SeenRequestAtFloor(alv_list.MyIP, a.Floor, a.Button)

			// if a.Button == 2 {
			// 	fsm.Fsm_onRequestButtonPress(&elev, &wld_view, alv_list.MyIP, &tmr, a.Floor, a.Button)
			// } else {
			// 	fmt.Println("Step 1")
			// 	wld_view.SetHallRequestAtFloor(a.Floor, int(a.Button))
			// 	go func() {
			// 		wld_updated <- true
			// 	} ()
			// }

			// Press of button shall update my worldview which will then propagate out and be published that new info has been found.
			// 		But we must seperate between cab and hall buttons since cab calls can only be handled by itself.
			// We must then have an own function for reading in the world view and update the requests matrix of the elevator.
			// Must find out how to make elevator only considers its requests matrix and state to make decisions.

			// When we change state we also need to update the world view, this should probably happen when we arrive at new floor and update state.

			// fmt.Printf("%+v\n", a)
			// driver.SetButtonLamp(a.Button, a.Floor, true)
			// fsm.Fsm_onRequestButtonPress(&elev, &tmr, a.Floor, a.Button)
			// fmt.Printf("Request floor: %d\n", a.Floor)

		case a := <-drv_floors:
			// fmt.Printf("This floor polled: %d\n", a)
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

		case <-ord_updated:
			// fmt.Println("Step 5")
			go func() { 
				world_view.SetAllLights(lgt_array)
				for floor, buttons := range wld_view.GetMyAssignedOrders(alv_list.MyIP) {
					for button, value := range buttons {
						if value {
							fsm.Fsm_onRequestButtonPress(&elev, &wld_view, alv_list.MyIP, &tmr, floor, driver.ButtonType(button))
							
						} else {
							elev.Request[floor][button] = 0
						}
					}
				}
				for floor,value := range wld_view.GetMyCabRequests(alv_list.MyIP) {
					if value {
						fsm.Fsm_onRequestButtonPress(&elev, &wld_view, alv_list.MyIP, &tmr, floor, driver.BT_Cab)
					} else {
						elev.Request[floor][driver.BT_Cab] = 0
					}
				}
			} ()
			
		case <-wld_updated:
			
			// fmt.Println("Step 3")

			go func() {	
				if alv_list.AmIMaster() {
					order_assigner.AssignOrders(&wld_view, &alv_list)
					ord_updated <- true
				}
			} ()
		}
	}
}
