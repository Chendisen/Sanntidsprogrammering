package world_view

import (
	"Sanntid/cyclic_counter"
	"Sanntid/driver"
)

// TODO: Have structs that is similar to the ones we send in messages
// 			They will correspond to our world view and act as a middleman 
// 			for fault checking messages before taking decisions. 
// 			Must therefore have functions that compares the received 
// 			messages and the ones of our world view.  

type AliveList struct {
	Version 	cyclic_counter.Counter
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

func UpdateWorldView(newView WorldView, currentView *WorldView, senderIP string, myIP string, aliveList AliveList, c chan int){

	currentView.AddNewNodes(newView)
	(&newView).AddNewNodes(*currentView)

	for i, floor := range newView.HallRequests {
		for j, hallRequest := range floor {
			if cyclic_counter.ShouldUpdate(hallRequest, currentView.HallRequests[i][j]) {
				cyclic_counter.UpdateValue(&currentView.HallRequests[i][j], hallRequest.Value)
			}
		}
	}

	for IP, NodeState := range newView.States {
		if IP != myIP{
			if(cyclic_counter.ShouldUpdate(NodeState.Version, currentView.States[IP].Version)){
				currentView.States[IP] = NodeState
			}
		}
	}

	if senderIP == aliveList.Master {
		for i, floor := range newView.AssignedOrders[myIP]{
			for j, orderAssigned := range floor{
				if orderAssigned != currentView.AssignedOrders[myIP][i][j]{
					c <- 1
					break
				}
			}
		}
		currentView.AssignedOrders = newView.AssignedOrders
	}
}

func (wv *WorldView) MakeWorldView(myIP string){
	for i := 0; i < driver.N_FLOORS; i++ {
		wv.HallRequests = append(wv.HallRequests, [2]cyclic_counter.Counter{cyclic_counter.MakeCounter(cyclic_counter.MAX), cyclic_counter.MakeCounter(cyclic_counter.MAX)})
	}

	wv.States[myIP] = MakeElevatorState()
	wv.AssignedOrders[myIP] = make([][2]bool, driver.N_FLOORS)
}

func (wv *WorldView) GetMyAssignedOrders(myIP string) [][2]bool{
	return wv.AssignedOrders[myIP]
} 

func (wv *WorldView) GetMyCabOrders(myIP string) []bool{
	return wv.States[myIP].GetCabRequests()
}

