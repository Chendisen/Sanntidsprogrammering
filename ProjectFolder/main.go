package main

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/fsm"
	"Sanntid/network"
	"Sanntid/order_assigner"
	"Sanntid/process_pair"
	"Sanntid/timer"
	"Sanntid/watchdog"
	"Sanntid/world_view"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	const numFloors int = 4
	const watchdogTime float64 = 3

	var elev elevator.Elevator = elevator.Elevator_uninitialized()
	var tmr_door timer.Timer = timer.Timer_uninitialized()
	var tmr_watchdog timer.Timer = timer.Timer_uninitialized()
	var alv_list world_view.AliveList = world_view.MakeAliveList()

	idFlag := flag.Int("id", 1, "Specifies ID of terminal")
	flag.Parse()
	myID := *idFlag
	alv_list.MyIP = fmt.Sprintf("%d", myID)
	alv_list.Master = fmt.Sprintf("%d", myID)
	alv_list.NodesAlive[0] = fmt.Sprintf("%d", myID)

	var wld_view world_view.WorldView = world_view.MakeWorldView(alv_list.MyIP)
	var hrd_list world_view.HeardFromList = world_view.MakeHeardFromList(alv_list.MyIP)
	
	lgt_array := make([][3]bool, numFloors)

	drv_buttons := make(chan driver.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	ord_updated := make(chan bool, 10)
	wld_updated := make(chan bool, 10)

	elv_dead := make(chan bool)
	start_new := make(chan bool)



	go process_pair.ProcessPair(alv_list.MyIP, &wld_view, &tmr_door, start_new)

	for range start_new{
		break
	}

	path2,_ := os.Getwd()
	cmd := exec.Command("gnome-terminal", "--window", "--", "sh", "-c", "cd "+path2+" && go run main.go")

	fmt.Printf("Path: %s", path2)

	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start myself")
		panic(err)
	}

	serverPort := 15656 + myID 
	serverPortString := fmt.Sprintf("%d", serverPort)
	driver.Init("localhost:"+serverPortString, numFloors)

	// driver.Init("localhost:15657", numFloors)


	go driver.PollButtons(drv_buttons)
	go driver.PollFloorSensor(drv_floors)
	go driver.PollObstructionSwitch(drv_obstr)
	go driver.PollStopButton(drv_stop)
	go fsm.Fsm_checkTimeOut(&elev, &wld_view, alv_list.MyIP, &tmr_door, &tmr_watchdog)
	go network.StartCommunication(alv_list.MyIP, &wld_view, &alv_list, &hrd_list, &lgt_array, ord_updated, wld_updated)
	go watchdog.Watchdog(&tmr_watchdog, &elev, elv_dead)

	fsm.Fsm_onInitBetweenFloors(&elev, &wld_view, alv_list.MyIP)
	world_view.InitLights(&lgt_array, alv_list.MyIP, wld_view)
	tmr_watchdog.Timer_start(watchdogTime)
	ord_updated<-true

	for {
		select {
		case a := <-drv_buttons:

			fmt.Println("A button pressed")

			wld_view.SeenRequestAtFloor(alv_list.MyIP, a.Floor, a.Button)

		case a := <-drv_floors:

			tmr_watchdog.Timer_start(watchdogTime)
			fsm.Fsm_onFloorArrival(&elev, &wld_view, alv_list.MyIP, &tmr_door, a)
			// wld_view.PrintWorldView()

		case a := <-drv_obstr:
			fmt.Printf("DOOR OBSTRUCTED: %t\n", a)
			elev.DoorObstructed = a
			wld_view.SetMyAvailabilityStatus(alv_list.MyIP, !a)
			go func() {
				if alv_list.AmIMaster() {
					wld_updated<-true
				}
			}()
			fmt.Printf("MY AVAILABILITY IS NOW: %t\n", wld_view.States[alv_list.MyIP].Available)

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := driver.ButtonType(0); b < 3; b++ {
					driver.SetButtonLamp(b, f, false)
				}
			}

		case <-ord_updated:

			go func() { 
				world_view.SetAllLights(lgt_array)
				if wld_view.GetMyAvailabilityStatus(alv_list.MyIP){
					for floor, buttons := range wld_view.GetMyAssignedOrders(alv_list.MyIP) {
						for button, value := range buttons {
							if value {
								fsm.Fsm_onRequestButtonPress(&elev, &wld_view, alv_list.MyIP, &tmr_door, &tmr_watchdog, floor, driver.ButtonType(button))
								
							} else {
								elev.Request[floor][button] = 0
							}
						}
					}
					for floor,value := range wld_view.GetMyCabRequests(alv_list.MyIP) {
						if value {
							fsm.Fsm_onRequestButtonPress(&elev, &wld_view, alv_list.MyIP, &tmr_door, &tmr_watchdog, floor, driver.BT_Cab)
						} else {
							elev.Request[floor][driver.BT_Cab] = 0
						}
					}
				}
			} ()
			
		case <-wld_updated:

			go func() {	
				if alv_list.AmIMaster() {
					order_assigner.AssignOrders(&wld_view, &alv_list)
				}
				ord_updated <- true
			} ()

		case <-elv_dead:

			// fmt.Println("THE ELEVATOR IS DEAD")
			panic("ELEVATOR DEAD")

			/*wld_view.SetMyAvailabilityStatus(alv_list.MyIP, false)
			fsm.Fsm_onInitBetweenFloors(&elev, &wld_view, alv_list.MyIP)
			
			go func() {
				if alv_list.AmIMaster() {
					wld_updated<-true
				}
				for a:= range drv_floors{
					wld_view.SetMyAvailabilityStatus(alv_list.MyIP, true)
					tmr_watchdog.Timer_start(watchdogTime)
					fsm.Fsm_onFloorArrival(&elev, &wld_view, alv_list.MyIP, &tmr_door, a)
					break
				}
			}()*/
		}
	}
}
