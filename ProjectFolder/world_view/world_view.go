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

func MakeElevatorState(counterMax int) ElevatorState{
	return ElevatorState{Version: cyclic_counter.MakeCounter(counterMax), Behaviour: "idle", Floor: -1, Direction: "stop", CabRequests: make([]bool, driver.N_FLOORS)}
}

func UpdateWorldView(newView WorldView, currentView *WorldView, senderIP string, myIP string, aliveList AliveList){
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
		currentView.AssignedOrders = newView.AssignedOrders
	}
}

func (wv *WorldView) InitWorldView(){
	var hallRequests [][2]cyclic_counter.Counter

	for i := 0; i < driver.N_FLOORS; i++ {
		hallRequests[i][0] = cyclic_counter.MakeCounter(cyclic_counter.MAX)
		hallRequests[i][1] = cyclic_counter.MakeCounter(cyclic_counter.MAX)
	}
}

func (wv *WorldView) GetMyAssignedOrders(myIP string) [][2]bool{
	return wv.AssignedOrders[myIP]
} 

