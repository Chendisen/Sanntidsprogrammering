package world_view

import (
	"Sanntid/resources/driver"
	. "Sanntid/resources/update_request"
)

type ElevatorState struct {
	Behaviour   string        `json:"behaviour"`
	Floor       int           `json:"floor"`
	Direction   string        `json:"direction"`
	CabRequests []OrderStatus `json:"cabRequests"`
	Available   bool          `json:"Available"`
}


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

func (elevatorState *ElevatorState) FinishedCabRequestAtFloor(floor int) {
	elevatorState.CabRequests[floor] = Order_Finished
}

func (elevatorState *ElevatorState) SetAvailabilityStatus(availabilityStatus bool) {
	elevatorState.Available = availabilityStatus
}

func (elevatorState ElevatorState) GetAvailabilityStatus() bool {
	return elevatorState.Available
}

func (elevatorState *ElevatorState) UpdateElevatorState(elv_update chan UpdateRequest) {
	for request := range elv_update {
		switch request.Type {
		case SetBehaviour:
			elevatorState.SetBehaviour(request.Value.(string))
		case SetFloor:
			elevatorState.SetFloor(request.Value.(int))
		case SetDirection:
			elevatorState.SetDirection(request.Value.(string))
		case SeenRequestAtFloor:
			elevatorState.SeenCabRequestAtFloor(request.Value.(int))
		case FinishedRequestAtFloor:
			elevatorState.FinishedCabRequestAtFloor(request.Value.(int))
		case SetMyAvailabilityStatus:
			elevatorState.SetAvailabilityStatus(request.Value.(bool))
		}
	}
}
