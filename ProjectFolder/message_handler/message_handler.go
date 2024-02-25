package message_handler

import(
	"encoding/json"
	"Sanntid/elevator"
	"Sanntid/cyclic_counter"
)

type ElevatorState struct {
    Behaviour   string      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   string      `json:"direction"`
    CabRequests []bool      `json:"cabRequests"`
}

type HRAInput struct {
    HallRequests    [][2]bool                   `json:"hallRequests"`
    States          map[string]ElevatorState    `json:"states"`
}

type StatesAndRequests struct {
    OrderAssignerInput	HRAInput				  	`json:"orderAssignerInput"`
	AssignedRequests	map[string][][2]bool		`json:"assignedRequests"`
}


type ElevatorStateMessage struct {
    Behaviour   string      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   string      `json:"direction"`
    CabRequests []int       `json:"cabRequests"`
	// TODO: Change data type of CabRequests to CyclicCounter
}

type HRAInputMessage struct {
    HallRequests    [][2]cyclic_counter.Counter			`json:"hallRequests"`
    States          map[string]ElevatorStateMessage     `json:"states"`
}

type StatesAndRequestsMessage struct {
    OrderAssignerInput	HRAInput				  					`json:"orderAssignerInput"`
	AssignedRequests	map[string][][2]cyclic_counter.Counter		`json:"assignedRequests"`
}

type StandardMessage struct {
    IPAddress		    string                   	`json:"IPAddress"`
  	StatesRequests		StatesAndRequestsMessage    `json:"statesRequests"`
}


func SetElevatorStateMessage(elevState elevator.Elevator) ElevatorStateMessage {
	var cabRequests []int
	for floor, requests := range(elevState.Request) {
		cabRequests = append(cabRequests, requests[2])
	}
	return ElevatorStateMessage{
		Behaviour:		elevState.Behaviour
		Floor:			elevState.Floor
		Direction:		elevState.Dirn
		CabRequests:	cabRequests
	}
}

// TODO: Should probably have functions that take in struct from wordlview
// 			and make the whole message, such that it can be called from main
// 			and then be given the correct variables. 

func GetSenderIP(message StandardMessage) string {
	return message.IPAddress
}

func GetAssignedRequests(message StandardMessage) map[string][][2]cyclic_counter.Counter{
	return message.StatesRequests.AssignedRequests
}

func GetStates(message StandardMessage) map[string]ElevatorStateMessage {
	return message.StatesRequests.OrderAssignerInput.States
}

func GetHallRequests(message StandardMessage) [][2]cyclic_counter.Counter {
	return message.StatesRequests.OrderAssignerInput.HallRequests
}

