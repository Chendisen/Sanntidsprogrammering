package world_view

import (
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/network/localip"
	"Sanntid/network/peers"
	"fmt"
	//"os"
)

// TODO: Have structs that is similar to the ones we send in messages
// 			They will correspond to our world view and act as a middleman
// 			for fault checking messages before taking decisions.
// 			Must therefore have functions that compares the received
// 			messages and the ones of our world view.

type AliveList struct {
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

func (os OrderStatus) ToBool() bool {
	return os == Order_Confirmed || os == Order_Finished
}

//ElevatorState functions

func MakeElevatorState() *ElevatorState {
	newElevator := new(ElevatorState)
	*newElevator = ElevatorState{Behaviour: "idle", Floor: -1, Direction: "stop", CabRequests: make([]OrderStatus, driver.N_FLOORS), Available: true}
	return newElevator
}

func (es ElevatorState) GetCabRequests() []bool {
	cabRequests := make([]bool, driver.N_FLOORS)
	for i, val := range es.CabRequests {
		cabRequests[i] = val.ToBool()
	}
	return cabRequests
}

func (es *ElevatorState) SetBehaviour(b string) {
	es.Behaviour = b
}

func (es *ElevatorState) SetFloor(f int) {
	es.Floor = f
}

func (es *ElevatorState) SetDirection(d string) {
	es.Direction = d
}

func (es *ElevatorState) SeenCabRequestAtFloor(f int) {
	es.CabRequests[f] = Order_Unconfirmed
}

func (es *ElevatorState) SetCabRequestAtFloor(f int) {
	es.CabRequests[f] = Order_Confirmed
}

func (es *ElevatorState) FinishedCabRequestAtFloor(f int) {
	es.CabRequests[f] = Order_Finished
}

func (es *ElevatorState) ClearCabRequestAtFloor(f int) {
	es.CabRequests[f] = Order_Empty
}

func (es *ElevatorState) SetAvailabilityStatus(availability_status bool) {
	es.Available = availability_status
}

func (es *ElevatorState) GetAvailabilityStatus() bool {
	return es.Available
}

//WordlView functions

//Requests

func (wv *WorldView) SetBehaviour(myIP string, eb elevator.ElevatorBehaviour) {
	wv.States[myIP].SetBehaviour(elevator.Eb_toString(eb))
}

func (wv *WorldView) SetFloor(myIP string, f int) {
	wv.States[myIP].SetFloor(f)
}

func (wv *WorldView) SetDirection(myIP string, md driver.MotorDirection) {
	wv.States[myIP].SetDirection(driver.Driver_dirn_toString(md))
}

func (wv *WorldView) SeenRequestAtFloor(myIP string, f int, b driver.ButtonType) {
	if b == driver.BT_Cab {
		if wv.States[myIP].CabRequests[f] == Order_Empty {
			wv.States[myIP].SeenCabRequestAtFloor(f)
		}
	} else {
		if wv.HallRequests[f][b] == Order_Empty {
			wv.HallRequests[f][b] = Order_Unconfirmed
		}
	}
}

/*func (wv *WorldView) SetRequestAtFloor(myIP string, f int, b driver.ButtonType) {
	if b == driver.BT_Cab {
		if wv.States[myIP].CabRequests[f] != Order_Empty{
			wv.States[myIP].FinishedCabRequestAtFloor(f)
		}
	} else {
		if wv.HallRequests[f][b] != Order_Empty{
			wv.HallRequests[f][b] = Order_Finished
		}
	}
}*/

func (wv *WorldView) FinishedRequestAtFloor(myIP string, f int, b driver.ButtonType) {
	if b == driver.BT_Cab {
		if wv.States[myIP].CabRequests[f] != Order_Empty {
			wv.States[myIP].FinishedCabRequestAtFloor(f)
		}
	} else {
		if wv.HallRequests[f][b] != Order_Empty {
			wv.HallRequests[f][b] = Order_Finished
		}
	}
}

/*func (wv *WorldView) ClearRequestAtFloor(myIP string, f int, b driver.ButtonType) {
	if b == driver.BT_Cab {
		wv.States[myIP].ClearCabRequestAtFloor(f)
	} else {
		wv.HallRequests[f][b] = Order_Empty
	}
}*/

func (wv WorldView) GetHallRequests() [][2]bool {
	var hall_requests [][2]bool = make([][2]bool, len(wv.HallRequests))
	for floor, buttons := range wv.HallRequests {
		for button, value := range buttons {
			hall_requests[floor][button] = value.ToBool()
		}
	}
	return hall_requests
}

func (wv *WorldView) GetMyAssignedOrders(myIP string) [][2]bool {
	return wv.AssignedOrders[myIP]
}

func (wv *WorldView) GetMyCabRequests(myIP string) []bool {
	return wv.States[myIP].GetCabRequests()
}

func (wv *WorldView) SetMyAvailabilityStatus(myIP string, availability_status bool) {
	wv.States[myIP].SetAvailabilityStatus(availability_status)
}

func (wv *WorldView) GetMyAvailabilityStatus(myIP string) bool {
	return wv.States[myIP].GetAvailabilityStatus()
}


//Nodes

func (wv *WorldView) ShouldAddNode(IP string) bool {
	if _, isPresent := wv.States[IP]; !isPresent {
		return true
	} else {
		return false
	}
}

func (wv *WorldView) AddNodeToWorldView(IP string) {
	wv.States[IP] = MakeElevatorState()
	wv.AssignedOrders[IP] = make([][2]bool, driver.N_FLOORS)
}

func (wv *WorldView) AddNewNodes(newView WorldView) {
	for IP := range newView.States {
		if wv.ShouldAddNode(IP) {
			wv.AddNodeToWorldView(IP)
		}
	}
}

//Updates

func (currentView *WorldView) UpdateWorldView(newView WorldView, senderIP string, sendTime string, myIP string, al AliveList, hfl *HeardFromList, lightArray *[][3]bool, ord_updated chan<- bool, wld_updated chan<- bool) {

	if senderIP == myIP {
		if !al.AmIMaster() {
			return
		}
	}

	currentView.AddNewNodes(newView)
	(&newView).AddNewNodes(*currentView)

	var wld_updated_flag bool = false
	var ord_updated_flag bool = false

	// fmt.Println("New view:")
	// newView.PrintWorldView()
	// fmt.Println("\nCurrent view:")
	// currentView.PrintWorldView()

	for f, floor := range newView.HallRequests {
		for b, buttonStatus := range floor {
			UpdateSynchronisedRequests(&currentView.HallRequests[f][b], buttonStatus, hfl, al, lightArray, f, b, senderIP, &wld_updated_flag, &ord_updated_flag, "")
		}
	}

	for IP, state := range newView.States {
		for f, floorStatus := range state.CabRequests {
			UpdateSynchronisedRequests(&currentView.States[IP].CabRequests[f], floorStatus, hfl, al, lightArray, f, driver.BT_Cab, senderIP, &wld_updated_flag, &ord_updated_flag, IP)
		}
	}

	if sendTime > currentView.LastHeard[senderIP] {
		currentView.States[senderIP].Behaviour = newView.States[senderIP].Behaviour
		currentView.States[senderIP].Direction = newView.States[senderIP].Direction
		currentView.States[senderIP].Floor = newView.States[senderIP].Floor
		currentView.States[senderIP].Available = newView.States[senderIP].Available
	}

	if wld_updated_flag {
		wld_updated <- true
	} else if ord_updated_flag {
		currentView.AssignedOrders = newView.AssignedOrders
		ord_updated <- true

	}

}

func MakeWorldView(myIP string) WorldView {
	var wv WorldView = WorldView{States: make(map[string]*ElevatorState), AssignedOrders: make(map[string][][2]bool), LastHeard: make(map[string]string)}

	for i := 0; i < driver.N_FLOORS; i++ {
		wv.HallRequests = append(wv.HallRequests, [2]OrderStatus{Order_Empty, Order_Empty})
	}

	wv.States[myIP] = MakeElevatorState()
	wv.AssignedOrders[myIP] = make([][2]bool, driver.N_FLOORS)

	return wv
}

func (wv WorldView) PrintWorldView() {
	for IP, states := range wv.States {

		fmt.Printf("State of %s: \n", IP)
		fmt.Printf("		Floor: %d\n", states.Floor)
		fmt.Printf("	Behaviour: %s\n", states.Behaviour)
		fmt.Printf("	Direction: %s\n", states.Direction)
		fmt.Println("")

	}

	/*fmt.Println("Hall requests: ")
	for f, floor := range wv.HallRequests {
		fmt.Printf("Floor: %d\n", f)
		for b, buttonStatus := range floor {
			fmt.Printf("	Button: %d, Status: %d\n", b, buttonStatus)
		}
	}

	fmt.Println("Cab requests: ")
	for IP,state := range wv.States {
		fmt.Printf("	Elevator: %s\n", IP)
		for f,buttonStatus := range state.CabRequests {
			fmt.Printf("		Floor: %d, Status: %d\n", f, buttonStatus)
		}
		fmt.Println("")
	}*/

}

//AliveList funcitons

func MakeAliveList() AliveList {
	myIP, _ := localip.LocalIP()
	nodesAlive := make([]string, 1)
	nodesAlive[0] = myIP
	//fmt.Printf("Length of nodesAlive: %d\n", len(nodesAlive))
	//myIP := os.Getpid()
	return AliveList{MyIP: myIP, NodesAlive: nodesAlive, Master: myIP}
}

func (al AliveList) AmIMaster() bool {
	if al.Master == al.MyIP {
		return true
	} else {
		return false
	}
}

func (al *AliveList) ShouldUpdateList(p peers.PeerUpdate) bool {
	if len(p.Lost) != 0 {
		return true
	} else if len(p.New) != 0 {
		return true
	} else {
		return false
	}
}

func (al *AliveList) ShouldUpdateMaster(p peers.PeerUpdate) (bool, string) {
	var shouldUpdate bool = false
	var newMaster string = ""
	if len(p.Lost) != 0 {
		for _, lostNode := range p.Lost {
			if lostNode == al.Master {
				shouldUpdate = true
				for _, candidate := range p.Peers {
					if candidate > newMaster {
						newMaster = candidate
					}
				}
				return shouldUpdate, newMaster
			}
		}
	} else if p.New > al.Master {
		newMaster = p.New
		shouldUpdate = true
		return shouldUpdate, newMaster
	}
	return shouldUpdate, newMaster
}

func (al *AliveList) UpdateMaster(newMaster string) {
	al.Master = newMaster
}

func (al *AliveList) UpdateAliveList(p peers.PeerUpdate) {
	al.NodesAlive = p.Peers
	shouldUpdateMaster, newMaster := al.ShouldUpdateMaster(p)

	if shouldUpdateMaster {
		al.UpdateMaster(newMaster)
	}
}

func (al AliveList) Print() {
	fmt.Printf("Current alive nodes: \n")
	for _, IP := range al.NodesAlive {
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

func (hfl HeardFromList) ShouldResetAtFloorButton(f int, b int, al AliveList) bool {
	var count int = 0
	for _, buttonArray := range hfl.HeardFrom {
		if buttonArray[f][b] {
			count++
		}
	}
	return count == len(al.NodesAlive)
}

func (hfl HeardFromList) ShouldAddNode(ip string) bool {
	var check bool = true
	for IP := range hfl.HeardFrom {
		if IP == ip {
			check = false
			return check
		}
	}
	return check
}

func (hfl *HeardFromList) SetHeardFrom(msgIP string, f int, b int) {
	hfl.HeardFrom[msgIP][f][b] = true
}

func (hfl *HeardFromList) GetHeardFrom(msgIP string, f int, b int) bool {
	return hfl.HeardFrom[msgIP][f][b]
}

func (hfl *HeardFromList) CheckHeardFromAll(alv_list AliveList, f int, b int) bool {
	var heard_from_all bool = true
	for _, alv_nodes := range alv_list.NodesAlive {
		heard_from_all = heard_from_all && hfl.HeardFrom[alv_nodes][f][b]
	}
	return heard_from_all
}

func (hfl *HeardFromList) ClearHeardFrom(f int, b int) {
	for _, hfl_buttons := range hfl.HeardFrom {
		hfl_buttons[f][b] = false
	}
}

func (hfl *HeardFromList) AddNodeToList(newIP string) {
	hfl.HeardFrom[newIP] = make([][3]bool, driver.N_FLOORS)
}

func (hfl HeardFromList) Print() {
	fmt.Println("We have heard from: ")
	for IP := range hfl.HeardFrom {
		fmt.Printf("	%s\n", IP)
	}
	fmt.Printf("")

	for IP, table := range hfl.HeardFrom {
		fmt.Printf("Elevator: %s \n", IP)
		for f, buttons := range table {
			fmt.Printf("	Floor: %d", f)
			for b := range buttons {
				fmt.Printf("	Button: %d", b)
			}
			fmt.Print("\n")
		}
	}
}

// Big switch case for update world view
func UpdateSynchronisedRequests(cur_req *OrderStatus, rcd_req OrderStatus, hfl *HeardFromList, alv_list AliveList, light_array *[][3]bool, f int, b int, rcd_IP string, wld_updated_flag *bool, ord_updated_flag *bool, cabIP string) {
	switch rcd_req {
	case Order_Empty: // No requests
		if *cur_req == Order_Finished {
			// TODO: Channel that turns off the lights
			if b == driver.BT_Cab && alv_list.MyIP == cabIP {
				(*light_array)[f][b] = false
			} else if b != driver.BT_Cab {
				(*light_array)[f][b] = false
			}
			*ord_updated_flag = true
			hfl.ClearHeardFrom(f, b)
			*cur_req = Order_Empty
			// fmt.Print("Case 1\n")
		}
	case Order_Unconfirmed: // Unconfirmed requests
		if *cur_req == Order_Empty || *cur_req == Order_Unconfirmed {
			*cur_req = Order_Unconfirmed
			hfl.SetHeardFrom(rcd_IP, f, b)
			if alv_list.AmIMaster() {
				if hfl.CheckHeardFromAll(alv_list, f, b) {
					// TODO: Channel for assigning orders
					// TODO: Channel for turning on the lights
					if b == driver.BT_Cab && alv_list.MyIP == cabIP {
						(*light_array)[f][b] = true
					} else if b != driver.BT_Cab {
						(*light_array)[f][b] = true
					}
					*wld_updated_flag = true
					hfl.ClearHeardFrom(f, b)
					*cur_req = Order_Confirmed
				}
			}
			// fmt.Print("Case 2\n")
		}
	case Order_Confirmed: // Confirmed requests
		if *cur_req == Order_Unconfirmed || *cur_req == Order_Empty{
			// TODO: Channel for updating assigned orders
			// TODO: Channel for turning on lights
			if b == driver.BT_Cab && alv_list.MyIP == cabIP {
				(*light_array)[f][b] = true
			} else if b != driver.BT_Cab {
				(*light_array)[f][b] = true
			}
			*ord_updated_flag = true
			hfl.ClearHeardFrom(f, b)
			*cur_req = Order_Confirmed
			// fmt.Print("Case 3\n")
		}
	case Order_Finished: // Finished requests
		if *cur_req == Order_Unconfirmed || *cur_req == Order_Confirmed || *cur_req == Order_Finished {
			*cur_req = Order_Finished
			hfl.SetHeardFrom(rcd_IP, f, b)
			if alv_list.AmIMaster() {
				if hfl.CheckHeardFromAll(alv_list, f, b) {
					// TODO: Channel for turning off lights
					if b == driver.BT_Cab && alv_list.MyIP == cabIP {
						(*light_array)[f][b] = false
					} else if b != driver.BT_Cab {
						(*light_array)[f][b] = false
					}
					*wld_updated_flag = true
					hfl.ClearHeardFrom(f, b)
					*cur_req = Order_Empty
				}
			}
			// fmt.Print("Case 4\n")
		}

	}
}

func SetAllLights(lightArray [][3]bool) {
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			//outputDevice.RequestButtonLight(floor, driver.ButtonType(btn), driver.IntToBool(es.Request[floor][btn]))
			driver.SetButtonLamp(driver.ButtonType(btn), floor, lightArray[floor][btn])
		}
	}
}

func InitLights(lightArray *[][3]bool, myIP string, wld_view WorldView){
	for floor, buttons := range wld_view.GetMyAssignedOrders(myIP) {
		for button, value := range buttons {
			(*lightArray)[floor][button] = value
		}
	}
	for floor,value := range wld_view.GetMyCabRequests(myIP) {
		(*lightArray)[floor][driver.BT_Cab] = value
	}
}

// TODO: Change design of fsm functions since we no longer set values of wld_view by ourselves.
// Only time we set it ourselves is when we receive an order and hallrequest value is set to one.
// And when we clear an order since we are finished and hallrequest value is set to three.
