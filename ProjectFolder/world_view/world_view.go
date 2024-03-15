package world_view

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	//"Sanntid/message_handler"
	"fmt"

	// "os"
	"time"
)


type WorldView struct {
	HallRequests   [][2]OrderStatus          `json:"hallRequests"`
	States         map[string]*ElevatorState `json:"states"`
	AssignedOrders map[string][][2]bool      `json:"assignedOrders"`
	LastHeard      map[string]string         `json:"lastHeard"`
}

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


//WordlView functions

func MakeWorldView(myIP string) WorldView {
	var worldView WorldView = WorldView{States: make(map[string]*ElevatorState), AssignedOrders: make(map[string][][2]bool), LastHeard: make(map[string]string)}

	for i := 0; i < driver.N_FLOORS; i++ {
		worldView.HallRequests = append(worldView.HallRequests, [2]OrderStatus{Order_Empty, Order_Empty})
	}

	worldView.States[myIP] = MakeElevatorState()
	worldView.AssignedOrders[myIP] = make([][2]bool, driver.N_FLOORS)

	return worldView
}

func (worldView *WorldView) SetBehaviour(myIP string, elevatorBehaviour elevator.ElevatorBehaviour) {
	worldView.States[myIP].SetBehaviour(elevator.ElevatorBehaviourToString(elevatorBehaviour))
}

func (worldView *WorldView) SetFloor(myIP string, floor int) {
	worldView.States[myIP].SetFloor(floor)
}

func (worldView *WorldView) SetDirection(myIP string, motorDirection driver.MotorDirection) {
	worldView.States[myIP].SetDirection(driver.DriverDirectionToString(motorDirection))
}

func (worldView *WorldView) SeenRequestAtFloor(myIP string, floor int, button driver.ButtonType) {
	if button == driver.BT_Cab {
		if worldView.States[myIP].CabRequests[floor] == Order_Empty {
			worldView.States[myIP].SeenCabRequestAtFloor(floor)
		}
	} else {
		if worldView.HallRequests[floor][button] == Order_Empty {
			worldView.HallRequests[floor][button] = Order_Unconfirmed
		}
	}
}

func (worldView *WorldView) FinishedRequestAtFloor(myIP string, floor int, button driver.ButtonType) {
	if button == driver.BT_Cab {
		if worldView.States[myIP].CabRequests[floor] != Order_Empty {
			worldView.States[myIP].FinishedCabRequestAtFloor(floor)
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

func (worldView WorldView) GetMyAssignedOrders(myIP string) [][2]bool {
	return worldView.AssignedOrders[myIP]
}

func (worldView WorldView) GetMyCabRequests(myIP string) []bool {
	return worldView.States[myIP].GetCabRequests()
}

func (worldView *WorldView) SetMyAvailabilityStatus(myIP string, availabilityStatus bool) {
	worldView.States[myIP].SetAvailabilityStatus(availabilityStatus)
}

func (worldView WorldView) GetMyAvailabilityStatus(myIP string) bool {
	return worldView.States[myIP].GetAvailabilityStatus()
}


//Nodes

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

//Updates

func (currentView *WorldView) UpdateWorldViewOnIncomingMessage(incomingMessage StandardMessage, myIP string, networkOverview NetworkOverview, heardFromList *HeardFromList, lightArray *LightArray, ord_updated chan<- bool, wld_updated chan<- bool) {

	
	newView := incomingMessage.WorldView
	senderIP := incomingMessage.IPAddress
	sendTime := incomingMessage.SendTime
	

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
func UpdateSynchronisedRequests(cur_req *OrderStatus, rcd_req OrderStatus, heardFromList *HeardFromList, networkOverview NetworkOverview, lightArray *LightArray, floor int, button int, rcd_IP string, wld_updated_flag *bool, ord_updated_flag *bool, cabIP string) {
	switch rcd_req {
	case Order_Empty: // No requests
		if *cur_req == Order_Finished {
			// TODO: Channel that turns off the lights
			if button == driver.BT_Cab && networkOverview.MyIP == cabIP {
				(*lightArray)[floor][button] = false
			} else if button != driver.BT_Cab {
				(*lightArray)[floor][button] = false
			}
			*ord_updated_flag = true
			heardFromList.ClearHeardFrom(floor, button)
			*cur_req = Order_Empty
		}
	case Order_Unconfirmed: // Unconfirmed requests
		if *cur_req == Order_Empty || *cur_req == Order_Unconfirmed {
			*cur_req = Order_Unconfirmed
			heardFromList.SetHeardFrom(networkOverview, rcd_IP, floor, button)
			if networkOverview.AmIMaster() {
				if heardFromList.CheckHeardFromAll(networkOverview, floor, button) {
					// TODO: Channel for assigning orders
					// TODO: Channel for turning on the lights
					if button == driver.BT_Cab && networkOverview.MyIP == cabIP {
						(*lightArray)[floor][button] = true
					} else if button != driver.BT_Cab {
						(*lightArray)[floor][button] = true
					}
					*wld_updated_flag = true
					heardFromList.ClearHeardFrom(floor, button)
					*cur_req = Order_Confirmed
				}
			}
		}
	case Order_Confirmed: // Confirmed requests
		if *cur_req == Order_Unconfirmed || *cur_req == Order_Empty{
			// TODO: Channel for updating assigned orders
			// TODO: Channel for turning on lights
			if button == driver.BT_Cab && networkOverview.MyIP == cabIP {
				(*lightArray)[floor][button] = true
			} else if button != driver.BT_Cab {
				(*lightArray)[floor][button] = true
			}
			*ord_updated_flag = true
			heardFromList.ClearHeardFrom(floor, button)
			*cur_req = Order_Confirmed
		}
	case Order_Finished: // Finished requests
		if *cur_req == Order_Unconfirmed || *cur_req == Order_Confirmed || *cur_req == Order_Finished {
			*cur_req = Order_Finished
			heardFromList.SetHeardFrom(networkOverview, rcd_IP, floor, button)
			if networkOverview.AmIMaster() {
				if heardFromList.CheckHeardFromAll(networkOverview, floor, button) {
					// TODO: Channel for turning off lights
					if button == driver.BT_Cab && networkOverview.MyIP == cabIP {
						(*lightArray)[floor][button] = false
					} else if button != driver.BT_Cab {
						(*lightArray)[floor][button] = false
					}
					*wld_updated_flag = true
					heardFromList.ClearHeardFrom(floor, button)
					*cur_req = Order_Empty
				}
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

func (worldView *WorldView) UpdateWorldView(networkOverview *NetworkOverview, heardFromList *HeardFromList, lightArray *LightArray, ord_updated chan bool, wld_updated chan bool, setBehaviour chan elevator.ElevatorBehaviour, setFloor chan int, setDirection chan driver.MotorDirection, seenRequestAtFloor chan driver.ButtonEvent, finishedRequestAtFloor chan driver.ButtonEvent, setMyAvailabilityStatus chan bool, updateWorldViewOnIncomingMessage chan StandardMessage) {
	myIP := networkOverview.MyIP
	
	for{
		select {
		case behaviour := <-setBehaviour:
			worldView.SetBehaviour(myIP, behaviour)
		case newFloor := <-setFloor:
			worldView.SetFloor(myIP, newFloor)
		case newDirection := <-setDirection:
			worldView.SetDirection(myIP, newDirection)
		case newRequest := <-seenRequestAtFloor:
			worldView.SeenRequestAtFloor(myIP, newRequest.Floor, newRequest.Button)
		case finishedRequest := <-finishedRequestAtFloor:
			worldView.FinishedRequestAtFloor(myIP, finishedRequest.Floor, finishedRequest.Button)
		case availabilityStatus := <-setMyAvailabilityStatus:
			worldView.SetMyAvailabilityStatus(myIP, availabilityStatus)
		case incomingMessage := <-updateWorldViewOnIncomingMessage:
			worldView.UpdateWorldViewOnIncomingMessage(incomingMessage, myIP, *networkOverview, heardFromList, lightArray, ord_updated, wld_updated)

		}
	}
}



func (worldView *WorldView) UpdateWorldView2(updateRequest chan UpdateRequest, networkOverview *NetworkOverview, heardFromList *HeardFromList, lightArray *LightArray, ord_updated chan bool, wld_updated chan bool) {
	myIP := networkOverview.MyIP
	
	for request := range updateRequest {
		switch request.Type{
		case SetBehaviour:
			worldView.SetBehaviour(myIP, request.Value.(elevator.ElevatorBehaviour))
		case SetFloor:
			worldView.SetFloor(myIP, request.Value.(int))
		case SetDirection:
			worldView.SetDirection(myIP, request.Value.(driver.MotorDirection))
		case SeenRequestAtFloor:
			worldView.SeenRequestAtFloor(myIP, request.Value.(driver.ButtonEvent).Floor, request.Value.(driver.ButtonEvent).Button)
		case FinishedRequestAtFloor:
			worldView.FinishedRequestAtFloor(myIP, request.Value.(driver.ButtonEvent).Floor, request.Value.(driver.ButtonEvent).Button)
		case UpdateOnIncomingMessage:
			worldView.UpdateWorldViewOnIncomingMessage(request.Value.(StandardMessage), myIP, *networkOverview, heardFromList, lightArray, ord_updated, wld_updated)
		}
	}
}