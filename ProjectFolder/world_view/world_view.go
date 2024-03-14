package world_view

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/network/localip"
	"Sanntid/network/peers"
	"fmt"
	// "os"
	"time"
)

type NetworkOverview struct {
	MyIP       string
	NodesAlive []string
	Master     string
}

type OrderStatus int

const (
	Order_Empty OrderStatus = iota
	Order_Unconfirmed
	Order_Confirmed
	Order_Finished
)

type HeardFromList struct {
	HeardFrom map[string][][3]bool
}

type ElevatorState struct {
	Behaviour   	string        `json:"behaviour"`
	Floor       	int           `json:"floor"`
	Direction   	string        `json:"direction"`
	CabRequests 	[]OrderStatus `json:"cabRequests"`
	Available	bool 		  `json:"Available"`
}

type WorldView struct {
	HallRequests   [][2]OrderStatus          `json:"hallRequests"`
	States         map[string]*ElevatorState `json:"states"`
	AssignedOrders map[string][][2]bool      `json:"assignedOrders"`
	LastHeard      map[string]string         `json:"lastHeard"`
}

func (orderStatus OrderStatus) ToBool() bool {
	return orderStatus == Order_Confirmed || orderStatus == Order_Finished
}

//ElevatorState functions

func MakeElevatorState() *ElevatorState {
	newElevator := new(ElevatorState)
	*newElevator = ElevatorState{Behaviour: "idle", Floor: -1, Direction: "stop", CabRequests: make([]OrderStatus, driver.N_FLOORS), Available: true}
	return newElevator
}

func (elevatorState ElevatorState) GetCabRequests() []bool {
	cabRequests := make([]bool, driver.N_FLOORS)
	for i, val := range elevatorState.CabRequests {
		cabRequests[i] = val.ToBool()
	}
	return cabRequests
}

func (elevatorState *ElevatorState) SetBehaviour(behaviour string) {
	elevatorState.Behaviour = behaviour
}

func (elevatorState *ElevatorState) SetFloor(floor int) {
	elevatorState.Floor = floor
}

func (elevatorState *ElevatorState) SetDirection(direction string) {
	elevatorState.Direction = direction
}

func (elevatorState *ElevatorState) SeenCabRequestAtFloor(floor int) {
	elevatorState.CabRequests[floor] = Order_Unconfirmed
}

func (elevatorState *ElevatorState) SetCabRequestAtFloor(floor int) {
	elevatorState.CabRequests[floor] = Order_Confirmed
}

func (elevatorState *ElevatorState) FinishedCabRequestAtFloor(floor int) {
	elevatorState.CabRequests[floor] = Order_Finished
}

func (elevatorState *ElevatorState) ClearCabRequestAtFloor(floor int) {
	elevatorState.CabRequests[floor] = Order_Empty
}

func (elevatorState *ElevatorState) SetAvailabilityStatus(availabilityStatus bool) {
	elevatorState.Available = availabilityStatus
}

func (elevatorState *ElevatorState) GetAvailabilityStatus() bool {
	return elevatorState.Available
}

//WordlView functions

//Requests

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

func (worldView *WorldView) GetMyAssignedOrders(myIP string) [][2]bool {
	return worldView.AssignedOrders[myIP]
}

func (worldView *WorldView) GetMyCabRequests(myIP string) []bool {
	return worldView.States[myIP].GetCabRequests()
}

func (worldView *WorldView) SetMyAvailabilityStatus(myIP string, availabilityStatus bool) {
	worldView.States[myIP].SetAvailabilityStatus(availabilityStatus)
}

func (worldView *WorldView) GetMyAvailabilityStatus(myIP string) bool {
	return worldView.States[myIP].GetAvailabilityStatus()
}


//Nodes

func (worldView *WorldView) ShouldAddNode(IP string) bool {
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

func (currentView *WorldView) UpdateWorldView(newView WorldView, senderIP string, sendTime string, myIP string, networkOverview NetworkOverview, heardFromList *HeardFromList, lightArray *[][3]bool, ord_updated chan<- bool, wld_updated chan<- bool) {

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

func MakeWorldView(myIP string) WorldView {
	var worldView WorldView = WorldView{States: make(map[string]*ElevatorState), AssignedOrders: make(map[string][][2]bool), LastHeard: make(map[string]string)}

	for i := 0; i < driver.N_FLOORS; i++ {
		worldView.HallRequests = append(worldView.HallRequests, [2]OrderStatus{Order_Empty, Order_Empty})
	}

	worldView.States[myIP] = MakeElevatorState()
	worldView.AssignedOrders[myIP] = make([][2]bool, driver.N_FLOORS)

	return worldView
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

//NetworkOverview funcitons

func MakeNetworkOverview() NetworkOverview {
	myIP, _ := localip.LocalIP()
	//myIP := fmt.Sprintf("%d", os.Getpid())
	nodesAlive := make([]string, 1)
	nodesAlive[0] = myIP
	return NetworkOverview{MyIP: myIP, NodesAlive: nodesAlive, Master: myIP}
}

func (networkOverview NetworkOverview) GetMyIP() string {
	return networkOverview.MyIP
}

func (networkOverview NetworkOverview) AmIMaster() bool {
	if networkOverview.Master == networkOverview.MyIP {
		return true
	} else {
		return false
	}
}

func (networkOverview NetworkOverview) AmIAlive(myIP string) bool {
	var amIAlive bool = false
	for _, aliveNode := range networkOverview.NodesAlive {
		amIAlive = amIAlive || myIP == aliveNode
	}
	return amIAlive
}

func (networkOverview *NetworkOverview) ShouldUpdateList(p peers.PeerUpdate) bool {
	if len(p.Lost) != 0 {
		return true
	} else if len(p.New) != 0 {
		return true
	} else {
		return false
	}
}

func (networkOverview *NetworkOverview) ShouldUpdateMaster(p peers.PeerUpdate) (bool, string) {
	var shouldUpdate bool = false
	var newMaster string = ""
	if len(p.Lost) != 0 {
		for _, lostNode := range p.Lost {
			if lostNode == networkOverview.Master {
				shouldUpdate = true
				for _, candidate := range p.Peers {
					if candidate > newMaster {
						newMaster = candidate
					}
				}
				return shouldUpdate, newMaster
			}
		}
	} else if p.New > networkOverview.Master {
		newMaster = p.New
		shouldUpdate = true
		return shouldUpdate, newMaster
	}
	return shouldUpdate, newMaster
}

func (networkOverview *NetworkOverview) UpdateMaster(newMaster string) {
	networkOverview.Master = newMaster
}

func (networkOverview NetworkOverview) NetworkLost (p peers.PeerUpdate) bool {
	var networkGoing bool = false
	for _, aliveNode := range p.Peers {
		networkGoing = networkGoing || networkOverview.MyIP == aliveNode
	}
	return !networkGoing
}

func (networkOverview *NetworkOverview) UpdateNetworkOverview(p peers.PeerUpdate) {
	networkOverview.NodesAlive = p.Peers
	shouldUpdateMaster, newMaster := networkOverview.ShouldUpdateMaster(p)

	if shouldUpdateMaster {
		networkOverview.UpdateMaster(newMaster)
	}
}

func (networkOverview NetworkOverview) Print() {
	fmt.Printf("Current alive nodes: \n")
	for _, IP := range networkOverview.NodesAlive {
		fmt.Printf("A node	%s\n", IP)
	}
	fmt.Println("")
}

// ShouldResetList functions

func MakeHeardFromList(myIP string) HeardFromList {
	heardFromList := HeardFromList{HeardFrom: make(map[string][][3]bool)}
	heardFromList.HeardFrom[myIP] = make([][3]bool, driver.N_FLOORS)
	return heardFromList
}

func (heardFromList HeardFromList) ShouldResetAtFloorButton(floor int, button int, networkOverview NetworkOverview) bool {
	var count int = 0
	for _, buttonArray := range heardFromList.HeardFrom {
		if buttonArray[floor][button] {
			count++
		}
	}
	return count == len(networkOverview.NodesAlive)
}

func (heardFromList HeardFromList) ShouldAddNode(ip string) bool {
	var check bool = true
	for IP := range heardFromList.HeardFrom {
		if IP == ip {
			check = false
			return check
		}
	}
	return check
}

func (heardFromList *HeardFromList) SetHeardFrom(networkOverview NetworkOverview, msgIP string, floor int, button int) {
	for _,id := range networkOverview.NodesAlive {
		if id == msgIP{
			heardFromList.HeardFrom[msgIP][floor][button] = true
			return
		}
	}
}

func (heardFromList *HeardFromList) GetHeardFrom(networkOverview NetworkOverview, msgIP string, floor int, button int) bool {
	for _,id := range networkOverview.NodesAlive {
		if id == msgIP{
			return heardFromList.HeardFrom[msgIP][floor][button]
		}
	}
	return false
}

func (heardFromList *HeardFromList) CheckHeardFromAll(networkOverview NetworkOverview, floor int, button int) bool {
	var heard_from_all bool = true
	for _, alv_nodes := range networkOverview.NodesAlive {
		heard_from_all = heard_from_all && heardFromList.HeardFrom[alv_nodes][floor][button]
	}
	return heard_from_all
}

func (heardFromList *HeardFromList) ClearHeardFrom(floor int, button int) {
	for _, hfl_buttons := range heardFromList.HeardFrom {
		hfl_buttons[floor][button] = false
	}
}

func (heardFromList *HeardFromList) AddNodeToList(newIP string) {
	heardFromList.HeardFrom[newIP] = make([][3]bool, driver.N_FLOORS)
}

func (heardFromList HeardFromList) Print() {
	fmt.Println("We have heard from: ")
	for IP := range heardFromList.HeardFrom {
		fmt.Printf("	%s\n", IP)
	}
	fmt.Printf("")

	for IP, table := range heardFromList.HeardFrom {
		fmt.Printf("Elevator: %s \n", IP)
		for floor, buttons := range table {
			fmt.Printf("	Floor: %d", floor)
			for button := range buttons {
				fmt.Printf("	Button: %d", button)
			}
			fmt.Print("\n")
		}
	}
}

// Big switch case for update world view
func UpdateSynchronisedRequests(cur_req *OrderStatus, rcd_req OrderStatus, heardFromList *HeardFromList, networkOverview NetworkOverview, lightArray *[][3]bool, floor int, button int, rcd_IP string, wld_updated_flag *bool, ord_updated_flag *bool, cabIP string) {
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

func SetAllLights(lightArray [][3]bool) {
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			driver.SetButtonLamp(driver.ButtonType(btn), floor, lightArray[floor][btn])
		}
	}
}

func InitLights(lightArray *[][3]bool, myIP string, worldView WorldView){
	for floor, buttons := range worldView.GetHallRequests() {
		for button, value := range buttons {
			(*lightArray)[floor][button] = value
		}
	}
	for floor,value := range worldView.GetMyCabRequests(myIP) {
		(*lightArray)[floor][driver.BT_Cab] = value
	}
}
