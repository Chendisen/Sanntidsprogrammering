package main

import (
	"Sanntid/communication"
	"Sanntid/elevator"
	"Sanntid/elevator/stop_button"
	"Sanntid/order_assigner"
	"Sanntid/process_pair"
	"Sanntid/resources/driver"
	. "Sanntid/resources/update_request"
	"Sanntid/timer"
	"Sanntid/timer/door_open_timer"
	"Sanntid/timer/watchdog"
	"Sanntid/world_view"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	idFlag := flag.Int("id", 1, "Specifies ID of terminal")
	flag.Parse()
	myID := *idFlag

	var elev elevator.Elevator = elevator.Elevator_uninitialized()
	var timerDoor timer.Timer = timer.TimerUninitialized()
	var timerWatchdog timer.Timer = timer.TimerUninitialized()
	var networkOverview world_view.NetworkOverview = world_view.MakeNetworkOverviewWithIDFlag(fmt.Sprintf("%d", myID))
	var worldView world_view.WorldView = world_view.MakeWorldView(networkOverview.GetMyIP())
	var heardFromList world_view.HeardFromList = world_view.MakeHeardFromList(networkOverview.GetMyIP())
	var lightArray elevator.LightArray = elevator.MakeLightArray()

	drv_buttons := make(chan driver.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	ord_updated := make(chan bool, 10)
	wld_updated := make(chan bool, 10)
	elev_dead := make(chan bool)
	start_new := make(chan bool)
	msg_received := make(chan world_view.StandardMessage, 10)
	upd_request := make(chan UpdateRequest, 10)
	

	// PROCESS PAIRS

	go process_pair.ProcessPair(networkOverview.GetMyIP(), &worldView, &timerDoor, start_new)

	for range start_new {
		break
	}

	path, _ := os.Getwd()
	cmd := exec.Command("gnome-terminal", "--window", "--", "sh", "-c", "cd "+path+" && go run main.go -id "+fmt.Sprintf("%d",myID))
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start myself")
		panic(err)
	}


	// ELEVATORSERVER

	serverPort := 15656 + myID
	serverPortString := fmt.Sprintf("%d", serverPort)
	driver.Init("localhost:"+serverPortString, driver.N_FLOORS)


	// ELEVATORPROGRAM

	go driver.PollButtons(drv_buttons)
	go driver.PollFloorSensor(drv_floors)
	go driver.PollObstructionSwitch(drv_obstr)
	go driver.PollStopButton(drv_stop)
	go door_open_timer.CheckDoorOpenTimeout(&elev, networkOverview.GetMyIP(), &timerDoor, &timerWatchdog, upd_request)
	go communication.StartCommunication(&worldView, &networkOverview, msg_received, &heardFromList, ord_updated, wld_updated)
	go watchdog.CheckWatchdogTimeout(&timerWatchdog, &elev, elev_dead)
	go worldView.UpdateWorldView(upd_request, msg_received, &networkOverview, &heardFromList, &lightArray, ord_updated, wld_updated)

	elevator.Fsm_onInitBetweenFloors(&elev, networkOverview.MyIP, upd_request)
	elevator.Fsm_initAllOrders(ord_updated)
	lightArray.InitLights(worldView.GetHallRequests(), worldView.GetMyCabRequests(networkOverview.GetMyIP()))
	timerWatchdog.TimerStart(timer.WATCHDOG_TimeoutTime)

	for {
		select {
		case a := <-drv_buttons:
			upd_request <- GenerateUpdateRequest(SeenRequestAtFloor, a)

		case a := <-drv_floors:
			elevator.Fsm_onFloorArrival(&elev, networkOverview.MyIP, &timerDoor, a, upd_request)
			timerWatchdog.TimerStart(timer.WATCHDOG_TimeoutTime)

		case a := <-drv_obstr:
			elev.DoorObstructed = a
			upd_request <- GenerateUpdateRequest(SetMyAvailabilityStatus, !a)
			// Redistribute orders
			go func() {
				if networkOverview.AmIMaster() {
					wld_updated <- true
				}
			}()

		case <-drv_stop:
			for f := 0; f < driver.N_FLOORS; f++ {
				for b := driver.ButtonType(0); b < 3; b++ {
					driver.SetButtonLamp(b, f, false)
				}
			}
			stop_button.STOP()

		case <-ord_updated:
			elev.PrintRequest()
			go func() {
				lightArray.SetAllLights()
				if worldView.GetMyAvailabilityStatus(networkOverview.MyIP) {
					elevator.Fsm_setAssignedOrders(worldView.GetMyAssignedOrders(networkOverview.GetMyIP()), &elev, networkOverview.GetMyIP(), &timerDoor, &timerWatchdog, upd_request)
					elevator.Fsm_setCabOrders(worldView.GetMyCabRequests(networkOverview.GetMyIP()), &elev, networkOverview.GetMyIP(), &timerDoor, &timerWatchdog, upd_request)
				}
			}()

		case <-wld_updated:
			go func() {
				if networkOverview.AmIMaster() {
					order_assigner.AssignOrders(worldView, networkOverview, upd_request)
				}
				ord_updated <- true
			}()

		case <-elev_dead:
			panic("ELEVATOR DEAD")
		}
	}
}
