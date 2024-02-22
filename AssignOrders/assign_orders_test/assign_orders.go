package assignorders

// import (
// 	// "Sanntid/driver"
// 	"Sanntid/elevator"
// 	// "Sanntid/requests"
// 	"time"
// )

// type Req struct {
// 	active     bool
// 	assignedTo string
// }

// type State struct {
// 	id    string
// 	state elevator.Elevator
// 	time  time.Duration
// }

// func isUnassigned(r Req) bool {
// 	return r.active && r.assignedTo == ""
// }

// func filterReq(fn func(Req) bool, reqs [][]Req) [][]bool {
// 	var result [][]bool

// 	for _, req_list := range reqs {

// 		var mapped []bool
// 		for _, req := range req_list {
// 			mapped = append(mapped, fn(req))
// 		}

// 		result = append(result, mapped)
// 	}

// 	return result
// }

// func toReq(hallReqs [2][]bool) [2][]Req {
// 	var result [2][]Req

// 	for floor, reqsAtFloor := range hallReqs {
// 		for button, _ := range reqsAtFloor {
// 			result [floor][button] = Req {	active: 		hallReqs[floor][button],
// 											assignedTo: 	""}
// 		}
// 	}

// 	return result
// }

// // func withReqs(fn func(Req), s State, reqs [2][]Req) elevator.Elevator {
// // 	return s.state.
// // }