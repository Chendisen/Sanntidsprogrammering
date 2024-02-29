package world_view

import (
	"Sanntid/cyclic_counter"
)

// TODO: Have structs that is similar to the ones we send in messages
// 			They will correspond to our world view and act as a middleman 
// 			for fault checking messages before taking decisions. 
// 			Must therefore have functions that compares the received 
// 			messages and the ones of our world view.  


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
}