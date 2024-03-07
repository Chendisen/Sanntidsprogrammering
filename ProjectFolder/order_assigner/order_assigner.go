package order_assigner

import(
	"os/exec"
	"fmt"
	"encoding/json"
	"Sanntid/world_view"
    "runtime"

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


func AssignOrders(wld_view *world_view.WorldView, alv_list *world_view.AliveList) {

	hraExecutable := ""
    switch runtime.GOOS {
        case "linux":   hraExecutable  = "hall_request_assigner"
        case "windows": hraExecutable  = "hall_request_assigner.exe"
        default:        panic("OS not supported")
    }

    var states map[string]HRAElevState = make(map[string]HRAElevState)
    for _, alive_elevator := range alv_list.NodesAlive {
        for elevator, state := range wld_view.States{
            if alive_elevator == elevator  {
                states[elevator] =  HRAElevState{
                    Behavior: 		state.Behaviour,
                    Floor: 			state.Floor,
                    Direction:		state.Direction,
                    CabRequests: 	state.CabRequests,
                }
            }
        }
    }

	input := HRAInput{
		HallRequests: wld_view.GetHallRequests(),
		States: states,
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
	
	wld_view.AssignedOrders = *output
} 