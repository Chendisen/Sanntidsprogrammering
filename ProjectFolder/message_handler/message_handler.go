package message_handler

import(
	"encoding/json"
	"Sanntid/elevator"
	"Sanntid/cyclic_counter"
)


type AssignedRequestsMessage struct {
	AssignedOrders		map[string][][2]bool	 		`json:"assignedOrders"`
}

type StatesAndRequestsMessage struct {
    OrderAssignerInput	HRAInputMessage		  			`json:"orderAssignerInput"`
	AssignedOrders		AssignedRequestsMessage			`json:"assignedRequests"`
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

