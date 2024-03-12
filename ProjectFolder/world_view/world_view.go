package world_view

import (
	"Sanntid/cyclic_counter"
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/network/localip"
	"Sanntid/network/peers"
	"fmt"
)

// TODO: Have structs that is similar to the ones we send in messages
// 			They will correspond to our world view and act as a middleman
// 			for fault checking messages before taking decisions.
// 			Must therefore have functions that compares the received
// 			messages and the ones of our world view.

type Role int

const (
	Slave Role = iota
	Master
)

type AliveList struct {
	MyIP       string
	NodesAlive []string
	Master     string
}

const (
	Order_Empty   			OrderStatus = 0
	Order_Unconfirmed                   = 1
	Order_Confirmed                 	= 2
	Order_Finished						= 3
)

type HeardFromList struct {
	HeardFrom map[string][][3]bool
}

type ElevatorState struct {
	Version     cyclic_counter.Counter 		`json:"version"`
	Behaviour   string                 		`json:"behaviour"`
	Floor       int                    		`json:"floor"`
	Direction   string                 		`json:"direction"`
	CabRequests []cyclic_counter.Counter	`json:"cabRequests"`
}

type WorldView struct {
	HallRequests   [][2]cyclic_counter.Counter `json:"hallRequests"`
	States         map[string]*ElevatorState   `json:"states"`
	AssignedOrders map[string][][2]bool        `json:"assignedOrders"`
}

//ElevatorState functions

func MakeElevatorState() *ElevatorState {
	newElevator := new(ElevatorState)
	*newElevator = ElevatorState{Version: cyclic_counter.MakeCounter(50), Behaviour: "idle", Floor: -1, Direction: "stop", CabRequests: make([]cyclic_counter.Counter, driver.N_FLOORS)}
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
	cyclic_counter.Increment(&es.Version)
}

func (es *ElevatorState) SetFloor(f int) {
	es.Floor = f
	cyclic_counter.Increment(&es.Version)
}

func (es *ElevatorState) SetDirection(d string) {
	es.Direction = d
	cyclic_counter.Increment(&es.Version)
}

func (es *ElevatorState) SetCabRequestAtFloor(f int) {
	es.CabRequests[f] = true
	cyclic_counter.Increment(&es.Version)
}

func (es *ElevatorState) ClearCabRequestAtFloor(f int) {
	es.CabRequests[f] = false
	cyclic_counter.Increment(&es.Version)
}

//WordlView functions

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

func (wv *WorldView) SetBehaviour(myIP string, eb elevator.ElevatorBehaviour) {
	wv.States[myIP].SetBehaviour(elevator.Eb_toString(eb))
}

func (wv *WorldView) SetFloor(myIP string, f int) {
	wv.States[myIP].SetFloor(f)

}

func (wv *WorldView) SetDirection(myIP string, md driver.MotorDirection) {
	wv.States[myIP].SetDirection(driver.Driver_dirn_toString(md))
}

func (wv *WorldView) SetHallRequestAtFloor(f int, b int) {
	if wv.HallRequests[f][b].ToBool() {
		fmt.Println("Step 2, not set")
		return
	} else {
		fmt.Println("Step 2, set")
		cyclic_counter.Increment(&wv.HallRequests[f][b])
	}
}

func (wv *WorldView) ClearHallRequestAtFloor(f int, b int) {
	if wv.HallRequests[f][b].ToBool() {
		cyclic_counter.Increment(&wv.HallRequests[f][b])
	}
}

func (wv *WorldView) SetRequestAtFloor(myIP string, btn_floor int, btn_type int) {
	es := wv.States[myIP]

	if btn_type == 2 {
		es.SetCabRequestAtFloor(btn_floor)
	}
}

func (wv *WorldView) ClearRequestAtFloor(myIP string, btn_floor int, btn_type int) {
	es := wv.States[myIP]

	if btn_type == 2 {
		es.ClearCabRequestAtFloor(btn_floor)
	} else {
		wv.ClearHallRequestAtFloor(btn_floor, btn_type)
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

func (currentView *WorldView) UpdateWorldView(newView WorldView, senderIP string, myIP string, aliveList AliveList, ord_updated chan<- bool, wld_updated chan<- bool) {

	currentView.AddNewNodes(newView)
	(&newView).AddNewNodes(*currentView)

	fmt.Println("Current: ")
	currentView.PrintWorldView()
	fmt.Println("\n\nNew: ")
	newView.PrintWorldView()

	var hallRequestsUpdated bool = false
	for i, floor := range newView.HallRequests {
		for j, hallRequest := range floor {
			if cyclic_counter.ShouldUpdate(hallRequest, currentView.HallRequests[i][j]) {
				cyclic_counter.UpdateValue(&currentView.HallRequests[i][j], hallRequest.Value)
				hallRequestsUpdated = true
			}
		}
	}

	fmt.Printf("Is hall requests updated?: %t\n\n", hallRequestsUpdated)

	if hallRequestsUpdated {
		wld_updated <- true
	}

	for IP, NodeState := range newView.States {
		if IP != myIP {
			if cyclic_counter.ShouldUpdate(NodeState.Version, currentView.States[IP].Version) {
				*currentView.States[IP] = *NodeState
			}
		}
	}

	if senderIP == aliveList.Master {
		for i, floor := range newView.AssignedOrders[myIP] {
			for j, orderAssigned := range floor {
				if orderAssigned != currentView.AssignedOrders[myIP][i][j] {
					ord_updated <- true
					break
				}
			}
		}
		currentView.AssignedOrders = newView.AssignedOrders
	}

	/*if isUpdated {
		wld_updated <- true
	}*/
}

func MakeWorldView(myIP string) WorldView {
	var wv WorldView = WorldView{States: make(map[string]*ElevatorState), AssignedOrders: make(map[string][][2]bool)}

	for i := 0; i < driver.N_FLOORS; i++ {
		wv.HallRequests = append(wv.HallRequests, [2]cyclic_counter.Counter{cyclic_counter.MakeCounter(cyclic_counter.MAX), cyclic_counter.MakeCounter(cyclic_counter.MAX)})
	}

	wv.States[myIP] = MakeElevatorState()
	wv.AssignedOrders[myIP] = make([][2]bool, driver.N_FLOORS)

	return wv
}

func (wv *WorldView) GetMyAssignedOrders(myIP string) [][2]bool {
	return wv.AssignedOrders[myIP]
}

func (wv *WorldView) GetMyCabRequests(myIP string) []bool {
	return wv.States[myIP].GetCabRequests()
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

func (hfl *HeardFromList) CheckHeardFromAll(alv_list *AliveList, f int, b int) bool {
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
func UpdateSynchronisedRequests(cur_req *cyclic_counter.Counter, rcd_req *cyclic_counter.Counter, hfl *HeardFromList, alv_list *AliveList, f int, b int, rcd_IP string) {
	switch rcd_req.Value {
	case Order_Empty: // No requests
		if cur_req.Value == Order_Finished {
			// TODO: Channel that turns off the lights
			// TODO: Clear correct value on elevatorstate requests
			hfl.ClearHeardFrom(f, b)
			cur_req.Value = Order_Empty
		} 
	case Order_Unconfirmed: // Unconfirmed requests
		if cur_req.Value == Order_Empty || cur_req.Value == Order_Unconfirmed {
			cur_req.Value = Order_Unconfirmed
			hfl.SetHeardFrom(rcd_IP, f, b)
			if alv_list.AmIMaster() {
				if hfl.CheckHeardFromAll(alv_list, f, b) {
					// TODO: Channel for assigning orders
					// TODO: Channel for turning on the lights
					// TODO: Set correct value on elevatorstate requests
					hfl.ClearHeardFrom(f, b)
					cur_req.Value = Order_Confirmed
				}
			}
		}
	case Order_Confirmed: // Confirmed requests
		if cur_req.Value == Order_Unconfirmed {
			// TODO: Channel for updating assigned orders
			// TODO: Channel for turning on lights
			hfl.ClearHeardFrom(f, b)
			cur_req.Value = Order_Confirmed
		}
	case Order_Finished: // Finished requests
		if cur_req.Value == Order_Unconfirmed || cur_req.Value == Order_Confirmed || cur_req.Value == Order_Finished {
			cur_req.Value = Order_Finished
			hfl.SetHeardFrom(rcd_IP, f, b)
			if alv_list.AmIMaster() {
				if hfl.CheckHeardFromAll(alv_list, f, b) {
					// TODO: Channel for turning off lights
					// TODO: Clear correct value on elevatorstate requests
					hfl.ClearHeardFrom(f, b)
					cur_req.Value = Order_Empty
				}
			}
		}
		
	}
}

// TODO: Change design of fsm functions since we no longer set values of wld_view by ourselves. 
			// Only time we set it ourselves is when we receive an order and hallrequest value is set to one. 
			// And when we clear an order since we are finished and hallrequest value is set to three. 