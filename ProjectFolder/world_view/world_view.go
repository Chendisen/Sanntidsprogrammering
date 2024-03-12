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
	Order_Empty   			OrderStatus = iota
	Order_Unconfirmed                   
	Order_Confirmed                 	
	Order_Finished						
)

type HeardFromList struct {
	HeardFrom map[string][][3]bool
}

type ElevatorState struct {
	Behaviour   string                 		`json:"behaviour"`
	Floor       int                    		`json:"floor"`
	Direction   string                 		`json:"direction"`
	CabRequests []OrderStatus				`json:"cabRequests"`
}

type WorldView struct {
	HallRequests   [][2]OrderStatus 			`json:"hallRequests"`
	States         map[string]*ElevatorState   	`json:"states"`
	AssignedOrders map[string][][2]bool        	`json:"assignedOrders"`
	LastHeard 	   map[string]string			`json:"lastHeard"`
}

func (os OrderStatus) ToBool() bool{
	return os == Order_Confirmed || os == Order_Finished
 }
//ElevatorState functions

func MakeElevatorState() *ElevatorState {
	newElevator := new(ElevatorState)
	*newElevator = ElevatorState{Behaviour: "idle", Floor: -1, Direction: "stop", CabRequests: make([]OrderStatus, driver.N_FLOORS)}
	return newElevator
}

func (es ElevatorState) GetCabRequests() []bool {
	cabRequests := make([]bool, driver.N_FLOORS)
	for i,val := range es.CabRequests {
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

func (wv *WorldView) SeenRequestAtFloor(myIP string, f int, b driver.ButtonType){
	if b == driver.BT_Cab{
		wv.States[myIP].SeenCabRequestAtFloor(f)
	} else{
		wv.HallRequests[f][b] = Order_Unconfirmed
	}
}

func (wv *WorldView) SetRequestAtFloor(myIP string, f int, b driver.ButtonType) {
	if b == driver.BT_Cab{
		wv.States[myIP].SetCabRequestAtFloor(f)
	} else{
		wv.HallRequests[f][b] = Order_Confirmed
	}
}

func (wv *WorldView) FinishedRequestAtFloor(myIP string, f int, b driver.ButtonType) {
	if b == driver.BT_Cab{
		wv.States[myIP].FinishedCabRequestAtFloor(f)
	} else{
		wv.HallRequests[f][b] = Order_Finished
	}
}

func (wv *WorldView) ClearRequestAtFloor(myIP string, f int, b driver.ButtonType) {
	if b == driver.BT_Cab{
		wv.States[myIP].ClearCabRequestAtFloor(f)
	} else{
		wv.HallRequests[f][b] = Order_Empty
	}
}

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

func (currentView *WorldView) UpdateWorldView(newView WorldView, senderIP string, sendTime string, myIP string, al AliveList, hfl *HeardFromList, ord_updated chan<- bool, wld_updated chan<- bool, set_lights chan<- bool) {

	if senderIP == myIP {
		if !al.AmIMaster() {
			return
		}
	}

	// currentView.AddNewNodes(newView)
	// (&newView).AddNewNodes(*currentView)

	for f, floor := range newView.HallRequests {
		for b, buttonPressed := range floor {
			UpdateSynchronisedRequests(&currentView.HallRequests[f][b], buttonPressed, hfl, al, f, b, senderIP, wld_updated, set_lights)
		}
	}

	for IP,state := range newView.States {
		if IP != myIP{
			for f,floor := range state.CabRequests {
				UpdateSynchronisedRequests(&currentView.States[myIP].CabRequests[f], floor, hfl, al, f, int(driver.BT_Cab), senderIP, wld_updated, set_lights)
			}
		}
	}

	if sendTime > currentView.LastHeard[senderIP] && senderIP != myIP{
		fmt.Println("We are indeed updating the state")
		currentView.States[senderIP] = newView.States[senderIP]
	}

	if senderIP == al.Master {
		for i, floor := range newView.AssignedOrders[myIP] {
			for j, orderAssigned := range floor {
				if orderAssigned != currentView.AssignedOrders[myIP][i][j] {
					ord_updated <- true
					currentView.AssignedOrders = newView.AssignedOrders
					break
				}
			}
		}
	}

	/*if isUpdated {
		wld_updated <- true
	}*/
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
	/*fmt.Println("World View:")
	for IP,states := range wv.States {
		fmt.Printf("	Floor of %s: %d \n", IP, states.Floor)
	}*/

	fmt.Println("Assigned orders: ")
	for IP, table := range wv.AssignedOrders {
		fmt.Printf("Elevator: %s\n", IP)
		for floor, values := range table {
			fmt.Printf("	Floor: %d\n", floor)
			for button, isAssigned := range values {
				fmt.Printf("		Button: %d, %t\n", button, isAssigned)
			}
		}
	}

}

//AliveList funcitons

func MakeAliveList() AliveList {
	myIP, _ := localip.LocalIP()
	//myIP := os.Getpid()
	return AliveList{MyIP: myIP, Master: myIP}
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

// ShouldResetList functions

func MakeHeardFromList(myIP string) HeardFromList{
	heardFromList := HeardFromList{HeardFrom: make(map[string][][3]bool)}
	heardFromList.HeardFrom[myIP] = make([][3]bool, driver.N_FLOORS)

	return heardFromList
}

func (hfl HeardFromList) ShouldResetAtFloorButton(f int, b int, al AliveList) bool {
	var count int = 0
	for _,buttonArray := range hfl.HeardFrom{
		if buttonArray[f][b] {
			count++
		}
	}
	return count == len(al.NodesAlive)
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


// Big switch case for update world view
func UpdateSynchronisedRequests(cur_req *OrderStatus, rcd_req OrderStatus, hfl *HeardFromList, alv_list AliveList, f int, b int, rcd_IP string, wld_updated chan<- bool, set_lights chan<- bool) {
	switch rcd_req {
	case Order_Empty: // No requests
		if *cur_req == Order_Finished {
			// TODO: Channel that turns off the lights
			// TODO: Clear correct value on elevatorstate requests
			hfl.ClearHeardFrom(f, b)
			*cur_req = Order_Empty
		} 
	case Order_Unconfirmed: // Unconfirmed requests
		if *cur_req == Order_Empty || *cur_req == Order_Unconfirmed {
			*cur_req = Order_Unconfirmed
			hfl.SetHeardFrom(rcd_IP, f, b)
			if alv_list.AmIMaster() {
				if hfl.CheckHeardFromAll(alv_list, f, b) {
					// TODO: Channel for assigning orders
					wld_updated <- true
					// TODO: Channel for turning on the lights
					// TODO: Set correct value on elevatorstate requests
					hfl.ClearHeardFrom(f, b)
					*cur_req = Order_Confirmed
				}
			}
		}
	case Order_Confirmed: // Confirmed requests
		if *cur_req == Order_Unconfirmed {
			// TODO: Channel for updating assigned orders
			// TODO: Channel for turning on lights
			set_lights <- true

			hfl.ClearHeardFrom(f, b)
			*cur_req = Order_Confirmed
		}
	case Order_Finished: // Finished requests
		if *cur_req == Order_Unconfirmed || *cur_req == Order_Confirmed || *cur_req == Order_Finished {
			*cur_req = Order_Finished
			hfl.SetHeardFrom(rcd_IP, f, b)
			if alv_list.AmIMaster() {
				if hfl.CheckHeardFromAll(alv_list, f, b) {
					// TODO: Channel for turning off lights
					// TODO: Clear correct value on elevatorstate requests
					hfl.ClearHeardFrom(f, b)
					*cur_req = Order_Empty
				}
			}
		}
		
	}
}

// TODO: Change design of fsm functions since we no longer set values of wld_view by ourselves. 
			// Only time we set it ourselves is when we receive an order and hallrequest value is set to one. 
			// And when we clear an order since we are finished and hallrequest value is set to three. 

			