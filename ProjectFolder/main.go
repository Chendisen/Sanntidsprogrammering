package main

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/fsm"
	"Sanntid/communication"
	"Sanntid/order_assigner"
	"Sanntid/process_pair"
	"Sanntid/timer"
	"Sanntid/timer/watchdog"
	"Sanntid/timer/door_open_timer"
	"Sanntid/world_view"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	var elev elevator.Elevator = elevator.Elevator_uninitialized()
	var timerDoor timer.Timer = timer.Timer_uninitialized()
	var timerWatchdog timer.Timer = timer.Timer_uninitialized()
	var networkOverview world_view.NetworkOverview = world_view.MakeNetworkOverview()

	idFlag := flag.Int("id", 1, "Specifies ID of terminal")
	flag.Parse()
	myID := *idFlag
	networkOverview.MyIP = fmt.Sprintf("%d", myID)
	networkOverview.Master = fmt.Sprintf("%d", myID)
	networkOverview.NodesAlive[0] = fmt.Sprintf("%d", myID)

	var worldView world_view.WorldView = world_view.MakeWorldView(networkOverview.MyIP)
	var heardFromList world_view.HeardFromList = world_view.MakeHeardFromList(networkOverview.MyIP)
	
	lightArray := world_view.MakeLightArray()

	drv_buttons := make(chan driver.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	ord_updated := make(chan bool, 10)
	wld_updated := make(chan bool, 10)

	elev_dead := make(chan bool)
	start_new := make(chan bool)



	go process_pair.ProcessPair(networkOverview.MyIP, &worldView, &timerDoor, start_new)

	for range start_new{
		break
	}

	path,_ := os.Getwd()
	cmd := exec.Command("gnome-terminal", "--window", "--", "sh", "-c", "cd "+path+" && go run main.go")

	fmt.Printf("Path: %s", path)

	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start myself")
		panic(err)
	}

	serverPort := 15656 + myID 
	serverPortString := fmt.Sprintf("%d", serverPort)
	driver.Init("localhost:"+serverPortString, driver.N_FLOORS)

	// driver.Init("localhost:15657", driver.N_FLOORS)

	set_behaviour := make(chan elevator.ElevatorBehaviour)
	set_floor := make(chan int)
	set_direction := make(chan driver.MotorDirection)
	see_request := make(chan driver.ButtonEvent)
	fin_request := make(chan driver.ButtonEvent)
	set_availability := make(chan bool)
	upd_worldview := make(chan world_view.StandardMessage, 10)


	go driver.PollButtons(drv_buttons)
	go driver.PollFloorSensor(drv_floors)
	go driver.PollObstructionSwitch(drv_obstr)
	go driver.PollStopButton(drv_stop)
	go door_open_timer.CheckDoorOpenTimeout(&elev, &worldView, networkOverview.MyIP, &timerDoor, &timerWatchdog)
	go communication.StartCommunication(networkOverview.MyIP, &worldView, &networkOverview, upd_worldview, &heardFromList, &lightArray, ord_updated, wld_updated)
	go watchdog.CheckWatchdogTimeout(&timerWatchdog, &elev, elev_dead)
	go world_view.UpdateWorldView(&worldView, &networkOverview, &heardFromList, &lightArray, ord_updated, wld_updated, set_behaviour, set_floor, set_direction, see_request, fin_request, set_availability, upd_worldview)

	fsm.Fsm_onInitBetweenFloors(&elev, &worldView, networkOverview.MyIP)
	lightArray.InitLights(networkOverview.MyIP, worldView)
	timerWatchdog.Timer_start(timer.WATCHDOG_TimeoutTime)
	ord_updated<-true

	for {
		select {
		case a := <-drv_buttons:

			fmt.Println("A button pressed")

			// worldView.SeenRequestAtFloor(networkOverview.MyIP, a.Floor, a.Button)
			

		case a := <-drv_floors:

			timerWatchdog.Timer_start(timer.WATCHDOG_TimeoutTime)
			fsm.Fsm_onFloorArrival(&elev, &worldView, networkOverview.MyIP, &timerDoor, a)
			// worldView.PrintWorldView()

		case a := <-drv_obstr:
			fmt.Printf("DOOR OBSTRUCTED: %t\n", a)
			elev.DoorObstructed = a
			worldView.SetMyAvailabilityStatus(networkOverview.MyIP, !a)
			go func() {
				if networkOverview.AmIMaster() {
					wld_updated<-true
				}
			}()
			fmt.Printf("MY AVAILABILITY IS NOW: %t\n", worldView.States[networkOverview.MyIP].Available)

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < driver.N_FLOORS; f++ {
				for b := driver.ButtonType(0); b < 3; b++ {
					driver.SetButtonLamp(b, f, false)
				}
			}

		case <-ord_updated:

			go func() { 
				lightArray.SetAllLights()
				if worldView.GetMyAvailabilityStatus(networkOverview.MyIP){
					for floor, buttons := range worldView.GetMyAssignedOrders(networkOverview.MyIP) {
						for button, value := range buttons {
							if value {
								fsm.Fsm_onRequestButtonPress(&elev, &worldView, networkOverview.MyIP, &timerDoor, &timerWatchdog, floor, driver.ButtonType(button))
								
							} else {
								elev.Request[floor][button] = 0
							}
						}
					}
					for floor,value := range worldView.GetMyCabRequests(networkOverview.MyIP) {
						if value {
							fsm.Fsm_onRequestButtonPress(&elev, &worldView, networkOverview.MyIP, &timerDoor, &timerWatchdog, floor, driver.BT_Cab)
						} else {
							elev.Request[floor][driver.BT_Cab] = 0
						}
					}
				}
			} ()
			
		case <-wld_updated:

			go func() {	
				if networkOverview.AmIMaster() {
					order_assigner.AssignOrders(&worldView, &networkOverview)
				}
				ord_updated <- true
			} ()

		case <-elev_dead:

			// fmt.Println("THE ELEVATOR IS DEAD")
			panic("ELEVATOR DEAD")

			/*worldView.SetMyAvailabilityStatus(networkOverview.MyIP, false)
			fsm.Fsm_onInitBetweenFloors(&elev, &worldView, networkOverview.MyIP)
			
			go func() {
				if networkOverview.AmIMaster() {
					wld_updated<-true
				}
				for a:= range drv_floors{
					worldView.SetMyAvailabilityStatus(networkOverview.MyIP, true)
					timerWatchdog.Timer_start(timer.WATCHDOG_TimeoutTime)
					fsm.Fsm_onFloorArrival(&elev, &worldView, networkOverview.MyIP, &timerDoor, a)
					break
				}
			}()*/
		}
	}
}
