package order_assigner

import(
	"os/exec"
	"fmt"
	"encoding/json"
	"Sanntid/world_view"

)


type HRAElevState struct {
    Behavior    string      `json:"behaviour"`
    Floor       int         `json:"floor"` 
    Direction   string      `json:"direction"`
    CabRequests []bool      `json:"cabRequests"`
}

type HRAInput struct {
    HallRequests    [][2]bool                   `json:"hallRequests"`
    States          map[string]HRAElevState     `json:"states"`
}


func assign_orders(a_world_view *world_view.WorldView) {

	hraExecutable := ""
    switch runtime.GOOS {
        case "linux":   hraExecutable  = "hall_request_assigner"
        case "windows": hraExecutable  = "hall_request_assigner.exe"
        default:        panic("OS not supported")
    }

	input := HRAInput{
		HallRequests: a_world_view.HallRequests,
		States: map[string]HRAElevState{
			for elevator, state := range a_world_view.States{
				elevator: HRAElevState{
					Behavior: 		state.Behaviour,
					Floor: 			state.Floor,
					Direction:		state.Direction,
					CabRequests: 	state.CabRequests,
				},
			}
		}
	}

	jsonBytes, err := json.Marshal(input)
    if err != nil {
        fmt.Println("json.Marshal error: ", err)
        return
    }
    
    ret, err := exec.Command("./order_assigner"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
    if err != nil {
        fmt.Println("exec.Command error: ", err)
        fmt.Println(string(ret))
        return
    }

	output := new(map[string][][2]bool)
    err = json.Unmarshal(ret, &output)
    if err != nil {
        fmt.Println("json.Unmarshal error: ", err)
        return
    }
	
	a_world_view.AssignedOrders = output
} 