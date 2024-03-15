package world_view

import "Sanntid/driver"

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

func (elevatorState ElevatorState) GetAvailabilityStatus() bool {
	return elevatorState.Available
}
