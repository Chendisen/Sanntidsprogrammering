package world_view

import (
	"Sanntid/resources/driver"
	"Sanntid/elevator"
	. "Sanntid/resources/update_request"
	"fmt"
	"time"
)


type OrderStatus int

const (
	Order_Empty OrderStatus = iota
	Order_Unconfirmed
	Order_Confirmed
	Order_Finished
)

func (orderStatus OrderStatus) ToBool() bool {
	return orderStatus == Order_Confirmed || orderStatus == Order_Finished
}


type WorldView struct {
	HallRequests   [][2]OrderStatus          `json:"hallRequests"`
	States         map[string]*ElevatorState `json:"states"`
	AssignedOrders map[string][][2]bool      `json:"assignedOrders"`
	LastHeard      map[string]string         `json:"lastHeard"`
}

func MakeWorldView(myIP string) WorldView {
	var worldView WorldView = WorldView{States: make(map[string]*ElevatorState), AssignedOrders: make(map[string][][2]bool), LastHeard: make(map[string]string)}

	for i := 0; i < driver.N_FLOORS; i++ {
		worldView.HallRequests = append(worldView.HallRequests, [2]OrderStatus{Order_Empty, Order_Empty})
	}

	worldView.States[myIP] = MakeElevatorState()
	worldView.AssignedOrders[myIP] = make([][2]bool, driver.N_FLOORS)

	return worldView
}

func (worldView *WorldView) SetBehaviour(myIP string, elevatorBehaviour elevator.ElevatorBehaviour, elv_update chan UpdateRequest) {
	elv_update<- GenerateUpdateRequest(SetBehaviour, elevator.ElevatorBehaviourToString(elevatorBehaviour))
}

func (worldView *WorldView) SetFloor(myIP string, floor int, elv_update chan<- UpdateRequest) {
	elv_update<- GenerateUpdateRequest(SetFloor, floor)
}

func (worldView *WorldView) SetDirection(myIP string, motorDirection driver.MotorDirection, elv_update chan<- UpdateRequest) {
	elv_update<- GenerateUpdateRequest(SetDirection, driver.DriverDirectionToString(motorDirection)) 
}

func (worldView *WorldView) SeenRequestAtFloor(myIP string, floor int, button driver.ButtonType, elv_update chan<- UpdateRequest) {
	if button == driver.BT_Cab {
		if worldView.States[myIP].CabRequests[floor] == Order_Empty {
			elv_update<- GenerateUpdateRequest(SeenRequestAtFloor, floor)
		}
	} else {
		if worldView.HallRequests[floor][button] == Order_Empty {
			worldView.HallRequests[floor][button] = Order_Unconfirmed
		}
	}
}

func (worldView *WorldView) FinishedRequestAtFloor(myIP string, floor int, button driver.ButtonType, elv_update chan<- UpdateRequest) {
	if button == driver.BT_Cab {
		if worldView.States[myIP].CabRequests[floor] != Order_Empty {
			elv_update<- GenerateUpdateRequest(FinishedRequestAtFloor, floor)
		}
	} else {
		if worldView.HallRequests[floor][button] != Order_Empty {
			worldView.HallRequests[floor][button] = Order_Finished
		}
	}
}

func (worldView WorldView) GetHallRequests() [][2]bool {
	var hall_requests [][2]bool = make([][2]bool, len(worldView.HallRequests))
	for floor, buttons := range worldView.HallRequests {
		for button, value := range buttons {
			hall_requests[floor][button] = value.ToBool()
		}
	}
	return hall_requests
}

func (worldView *WorldView) SetAssignedOrders(assignedOrders map[string][][2]bool) {
	worldView.AssignedOrders = assignedOrders
}

func (worldView WorldView) GetMyAssignedOrders(myIP string) [][2]bool {
	return worldView.AssignedOrders[myIP]
}

func (worldView WorldView) GetMyCabRequests(myIP string) []bool {
	return worldView.States[myIP].GetCabRequests()
}

func (worldView *WorldView) SetMyAvailabilityStatus(myIP string, availabilityStatus bool, elv_update chan<- UpdateRequest) {
	elv_update<- GenerateUpdateRequest(SetMyAvailabilityStatus, availabilityStatus)
}

func (worldView WorldView) GetMyAvailabilityStatus(myIP string) bool {
	return worldView.States[myIP].GetAvailabilityStatus()
}


//ElevatorState Nodes

func (worldView WorldView) ShouldAddNode(IP string) bool {
	if _, isPresent := worldView.States[IP]; !isPresent {
		return true
	} else {
		return false
	}
}

func (worldView *WorldView) AddNodeToWorldView(IP string) {
	worldView.States[IP] = MakeElevatorState()
	worldView.AssignedOrders[IP] = make([][2]bool, driver.N_FLOORS)
}

func (worldView *WorldView) AddNewNodes(newView WorldView) {
	for IP := range newView.States {
		if worldView.ShouldAddNode(IP) {
			worldView.AddNodeToWorldView(IP)
		}
	}
}


// Updates

func (currentView *WorldView) UpdateWorldViewOnReceivedMessage(receivedMessage StandardMessage, myIP string, networkOverview NetworkOverview, heardFromList *HeardFromList, lightArray *elevator.LightArray, ord_updated chan<- bool, wld_updated chan<- bool) {

	newView := receivedMessage.WorldView
	senderIP := receivedMessage.IPAddress
	sendTime := receivedMessage.SendTime
	
	if senderIP == myIP {
		if !networkOverview.AmIMaster() {
			return
		}
	}

	currentView.AddNewNodes(newView)
	(&newView).AddNewNodes(*currentView)

	var wld_updated_flag bool = false
	var ord_updated_flag bool = false

	for floor, buttons := range newView.HallRequests {
		for button, buttonStatus := range buttons {
			UpdateSynchronisedRequests(&currentView.HallRequests[floor][button], buttonStatus, heardFromList, networkOverview, lightArray, floor, button, senderIP, &wld_updated_flag, &ord_updated_flag, "")
		}
	}

	for IP, state := range newView.States {
		for floor, floorStatus := range state.CabRequests {
			UpdateSynchronisedRequests(&currentView.States[IP].CabRequests[floor], floorStatus, heardFromList, networkOverview, lightArray, floor, driver.BT_Cab, senderIP, &wld_updated_flag, &ord_updated_flag, IP)
		}
	}

	if sendTime > currentView.LastHeard[senderIP] {
		currentView.States[senderIP].Behaviour = newView.States[senderIP].Behaviour
		currentView.States[senderIP].Direction = newView.States[senderIP].Direction
		currentView.States[senderIP].Floor = newView.States[senderIP].Floor
		if currentView.States[senderIP].Available != newView.States[senderIP].Available {
			wld_updated_flag = true
		}
		currentView.States[senderIP].Available = newView.States[senderIP].Available
	}

	if (senderIP == networkOverview.Master && sendTime > currentView.LastHeard[senderIP]) {
		currentView.AssignedOrders = newView.AssignedOrders
	}

	if wld_updated_flag {
		wld_updated <- true
	} else if ord_updated_flag {
		currentView.AssignedOrders = newView.AssignedOrders
		ord_updated <- true

	}

	currentView.LastHeard[senderIP] = time.Now().String()[11:19]
}

// Big switch case for update world view
func UpdateSynchronisedRequests(cur_req *OrderStatus, rcd_req OrderStatus, heardFromList *HeardFromList, networkOverview NetworkOverview, lightArray *elevator.LightArray, floor int, button int, rcd_IP string, wld_updated_flag *bool, ord_updated_flag *bool, cabIP string) {
	switch rcd_req {
	case Order_Empty: // No requests
		if *cur_req == Order_Finished {
			if button == driver.BT_Cab && networkOverview.MyIP == cabIP {
				lightArray.ClearElevatorLight(floor, button)
			} else if button != driver.BT_Cab {
				lightArray.ClearElevatorLight(floor, button)
			}
			*ord_updated_flag = true
			heardFromList.ClearHeardFrom(floor, button)
			*cur_req = Order_Empty
		}
	case Order_Unconfirmed: // Unconfirmed requests
		if *cur_req == Order_Empty || *cur_req == Order_Unconfirmed {
			*cur_req = Order_Unconfirmed
			heardFromList.SetHeardFrom(networkOverview, rcd_IP, floor, button)
			if networkOverview.AmIMaster() || heardFromList.CheckHeardFromAll(networkOverview, floor, button){
				if button == driver.BT_Cab && networkOverview.MyIP == cabIP {
					lightArray.SetElevatorLight(floor, button)
				} else if button != driver.BT_Cab {
					lightArray.SetElevatorLight(floor, button)
				}
				*wld_updated_flag = true
				heardFromList.ClearHeardFrom(floor, button)
				*cur_req = Order_Confirmed
			}
		}
	case Order_Confirmed: // Confirmed requests
		if *cur_req == Order_Unconfirmed || *cur_req == Order_Empty{
			if button == driver.BT_Cab && networkOverview.MyIP == cabIP {
				lightArray.SetElevatorLight(floor, button)
			} else if button != driver.BT_Cab {
				lightArray.SetElevatorLight(floor, button)
			}
			*ord_updated_flag = true
			heardFromList.ClearHeardFrom(floor, button)
			*cur_req = Order_Confirmed
		}
	case Order_Finished: // Finished requests
		if *cur_req == Order_Unconfirmed || *cur_req == Order_Confirmed || *cur_req == Order_Finished {
			*cur_req = Order_Finished
			heardFromList.SetHeardFrom(networkOverview, rcd_IP, floor, button)
			if networkOverview.AmIMaster() || heardFromList.CheckHeardFromAll(networkOverview, floor, button){
				if button == driver.BT_Cab && networkOverview.MyIP == cabIP {
					lightArray.ClearElevatorLight(floor, button)
				} else if button != driver.BT_Cab {
					lightArray.ClearElevatorLight(floor, button)
				}
				*wld_updated_flag = true
				heardFromList.ClearHeardFrom(floor, button)
				*cur_req = Order_Empty
			}
		}
	}
}

func (worldView WorldView) PrintWorldView() {
	/*for IP, states := range worldView.States {

		fmt.Printf("State of %s: \n", IP)
		fmt.Printf("		Floor: %d\n", states.Floor)
		fmt.Printf("	Behaviour: %s\n", states.Behaviour)
		fmt.Printf("	Direction: %s\n", states.Direction)
		fmt.Println("")

	}*/

	/*fmt.Println("Hall requests: ")
	for floor, floor := range worldView.HallRequests {
		fmt.Printf("Floor: %d\n", floor)
		for button, buttonStatus := range floor {
			fmt.Printf("	Button: %d, Status: %d\n", button, buttonStatus)
		}
	}

	fmt.Println("Cab requests: ")
	for IP,state := range worldView.States {
		fmt.Printf("	Elevator: %s\n", IP)
		for floor,buttonStatus := range state.CabRequests {
			fmt.Printf("		Floor: %d, Status: %d\n", floor, buttonStatus)
		}
		fmt.Println("")
	}*/

	fmt.Println("Assigned orders: ")
	for IP, orders := range worldView.AssignedOrders {
		fmt.Printf("	Elevator: %s\n", IP)
		for floor, buttons := range orders {
			fmt.Printf("		Floor: %d", floor)
			for button, value := range buttons {
				fmt.Printf("		Button: %d, Value: %t", button, value)
			}
			fmt.Print("\n")
		}
	}

}

func (worldView *WorldView) UpdateWorldView(upd_request chan UpdateRequest, msg_received chan StandardMessage, networkOverview *NetworkOverview, heardFromList *HeardFromList, lightArray *elevator.LightArray, ord_updated chan bool, wld_updated chan bool) {
	myIP := networkOverview.MyIP
	
	elv_update := make(chan UpdateRequest)
	go worldView.States[myIP].UpdateElevatorState(elv_update)

	for {
		select{
		case request := <-upd_request:	
			switch request.Type{
			case SetBehaviour:
				worldView.SetBehaviour(myIP, request.Value.(elevator.ElevatorBehaviour), elv_update)
			case SetFloor:
				worldView.SetFloor(myIP, request.Value.(int), elv_update)
			case SetDirection:
				worldView.SetDirection(myIP, (request.Value.(driver.MotorDirection)), elv_update)
			case SeenRequestAtFloor:
				worldView.SeenRequestAtFloor(myIP, request.Value.(driver.ButtonEvent).Floor, request.Value.(driver.ButtonEvent).Button, elv_update)
			case FinishedRequestAtFloor:
				worldView.FinishedRequestAtFloor(myIP, request.Value.(driver.ButtonEvent).Floor, request.Value.(driver.ButtonEvent).Button, elv_update)
			case SetAssignedOrders:
				worldView.SetAssignedOrders(request.Value.(map[string][][2]bool))
			case SetMyAvailabilityStatus:
				worldView.SetMyAvailabilityStatus(myIP, request.Value.(bool), elv_update)
			}
		case receivedMessage := <-msg_received:
			worldView.UpdateWorldViewOnReceivedMessage(receivedMessage, myIP, *networkOverview, heardFromList, lightArray, ord_updated, wld_updated)
		}
	}
}