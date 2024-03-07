package world_view

import (
	"Sanntid/cyclic_counter"
	"Sanntid/driver"
	"Sanntid/elevator"
	"Sanntid/network/localip"
	"Sanntid/network/peers"
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
	MyIP string
	NodesAlive 	[]string
	Master 		string
}

type ElevatorState struct {
	Version 		cyclic_counter.Counter  `json:"version"`
    Behaviour   	string      			`json:"behaviour"`
    Floor       	int         			`json:"floor"` 
    Direction   	string      			`json:"direction"`
    CabRequests 	[]bool      			`json:"cabRequests"`
}

type WorldView struct {
    HallRequests    [][2]cyclic_counter.Counter	`json:"hallRequests"`
    States          map[string]ElevatorState    `json:"states"`
	AssignedOrders  map[string][][2]bool		`json:"assignedOrders"` 			
}

//ElevatorState functions

func MakeElevatorState() ElevatorState{
	return ElevatorState{Version: cyclic_counter.MakeCounter(50), Behaviour: "idle", Floor: -1, Direction: "stop", CabRequests: make([]bool, driver.N_FLOORS)}
}

func (es ElevatorState) GetCabRequests() []bool{
	return es.CabRequests
} 

func (es *ElevatorState) SetBehaviour(b string){
	es.Behaviour = b
	cyclic_counter.Increment(&es.Version)
}

func (es *ElevatorState) SetFloor(f int){
	es.Floor = f
	cyclic_counter.Increment(&es.Version)
}

func (es *ElevatorState) SetDirection(d string){
	es.Direction = d
	cyclic_counter.Increment(&es.Version)
}

func (es *ElevatorState) SetCabRequestAtFloor(f int){
	es.CabRequests[f] = true
	cyclic_counter.Increment(&es.Version)
}

func (es *ElevatorState) ClearCabRequestAtFloor(f int){
	es.CabRequests[f] = false
	cyclic_counter.Increment(&es.Version)
}
//WordlView functions

func (wv *WorldView) ShouldAddNode(IP string) bool{
	if _,isPresent := wv.States[IP]; !isPresent {
		return true
	} else {
		return false
	}
}

func (wv *WorldView) AddNodeToWorldView(IP string){
	wv.States[IP] = MakeElevatorState()
	wv.AssignedOrders[IP] = make([][2]bool, driver.N_FLOORS)
}

func (wv *WorldView) AddNewNodes(newView WorldView){
	for IP := range newView.States{
		if wv.ShouldAddNode(IP) {
			wv.AddNodeToWorldView(IP)
		}
	}
}

func (wv *WorldView) SetBehaviour(myIP string, eb elevator.ElevatorBehaviour){
	es := wv.States[myIP]
	(&es).SetBehaviour()
}

func (wv *WorldView) SetFloor( myIP string, f int){
	es := wv.States[myIP]
	(&es).SetFloor(f)
}

func (wv *WorldView) SetDirection(myIP string, d string){
	es := wv.States[myIP]
	(&es).SetDirection(d)
}

func (wv *WorldView) SetHallRequestAtFloor(f int, b int){
	if(wv.HallRequests[f][b].ToBool()){
		return
	} else{
		cyclic_counter.Increment(&wv.HallRequests[f][b])
	}
}

func (wv *WorldView) ClearHallRequestAtFloor(f int, b int){
	if(wv.HallRequests[f][b].ToBool()){
		cyclic_counter.Increment(&wv.HallRequests[f][b])
	} 
}

func (wv *WorldView) SetRequestAtFloor(myIP string, btn_floor int, btn_type int) {
	es := wv.States[myIP]

	if btn_type == 2 {
		(&es).SetCabRequestAtFloor(btn_floor)
	} else {
		wv.SetHallRequestAtFloor(btn_floor, btn_type)
	}
}

func (wv *WorldView) ClearRequestAtFloor(myIP string, btn_floor int, btn_type int) {
	es := wv.States[myIP]

	if btn_type == 2 {
		(&es).ClearCabRequestAtFloor(btn_floor)
	} else {
		wv.ClearHallRequestAtFloor(btn_floor, btn_type)
	}
}

func (wv WorldView) GetHallRequests() [][2]bool {
	var hall_requests [][2]bool
	for floor, buttons := range wv.HallRequests {
		for button, value := range buttons {
			hall_requests[floor][button] = value.ToBool()
		}
	}
	return hall_requests
}

func (currentView *WorldView) UpdateWorldView(newView WorldView, senderIP string, myIP string, aliveList AliveList, ord_updated chan<- bool, wld_updated chan<- bool) bool{

	var isUpdated bool = false

	currentView.AddNewNodes(newView)
	(&newView).AddNewNodes(*currentView)

	for i, floor := range newView.HallRequests {
		for j, hallRequest := range floor {
			if cyclic_counter.ShouldUpdate(hallRequest, currentView.HallRequests[i][j]) {
				cyclic_counter.UpdateValue(&currentView.HallRequests[i][j], hallRequest.Value)
				isUpdated = true
			}
		}
	}

	for IP, NodeState := range newView.States {
		if IP != myIP{
			if(cyclic_counter.ShouldUpdate(NodeState.Version, currentView.States[IP].Version)){
				currentView.States[IP] = NodeState
				isUpdated = true
			}
		}
	}

	if senderIP == aliveList.Master {
		for i, floor := range newView.AssignedOrders[myIP]{
			for j, orderAssigned := range floor{
				if orderAssigned != currentView.AssignedOrders[myIP][i][j]{
					ord_updated <- true
					break
				}
			}
		}
		currentView.AssignedOrders = newView.AssignedOrders
		isUpdated = true
	}

	if isUpdated {
		wld_updated <- true
	}

	return isUpdated
}

func MakeWorldView(myIP string) WorldView{
	var wv WorldView = WorldView{States: make(map[string]ElevatorState), AssignedOrders: make(map[string][][2]bool)}
	
	for i := 0; i < driver.N_FLOORS; i++ {
		wv.HallRequests = append(wv.HallRequests, [2]cyclic_counter.Counter{cyclic_counter.MakeCounter(cyclic_counter.MAX), cyclic_counter.MakeCounter(cyclic_counter.MAX)})
	}

	wv.States[myIP] = MakeElevatorState()
	wv.AssignedOrders[myIP] = make([][2]bool, driver.N_FLOORS)

	return wv
}

func (wv *WorldView) GetMyAssignedOrders(myIP string) [][2]bool{
	return wv.AssignedOrders[myIP]
} 

func (wv *WorldView) GetMyCabRequests(myIP string) []bool{
	return wv.States[myIP].GetCabRequests()
}


//AliveList funcitons

func MakeAliveList() AliveList{
	myIP,_ := localip.LocalIP()
	return AliveList{MyIP: myIP, Master: myIP}
}

func (al AliveList) AmIMaster() bool {
	if al.Master == al.MyIP {
		return true
	} else {
		return false
	}
}

func (al *AliveList) ShouldUpdateList(p peers.PeerUpdate) bool{
	if len(p.Lost) != 0 {
		return true
	} else if len(p.New) != 0 {
		return true
	} else {
		return false
	}
}

func (al *AliveList) ShouldUpdateMaster(p peers.PeerUpdate) (bool, string){
	var shouldUpdate bool = false
	var newMaster string = ""
	if len(p.Lost) != 0{
		for _,lostNode := range p.Lost {
			if lostNode == al.Master{
				shouldUpdate = true
				for _,candidate := range p.Peers {
					if candidate > newMaster{
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

func (al *AliveList) UpdateMaster(newMaster string){
	al.Master = newMaster
}

func (al *AliveList) UpdateAliveList(p peers.PeerUpdate){
	al.NodesAlive = p.Peers
	shouldUpdateMaster, newMaster := al.ShouldUpdateMaster(p)

	if shouldUpdateMaster {
		al.UpdateMaster(newMaster)
	}
}

