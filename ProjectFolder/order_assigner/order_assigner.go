package order_assigner

import (
	"Sanntid/world_view"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func AssignOrders(wld_view *world_view.WorldView, alv_list *world_view.AliveList) {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	var states map[string]HRAElevState = make(map[string]HRAElevState)
	for _, alive_elevator := range alv_list.NodesAlive {
		if wld_view.States[alive_elevator].GetAvailabilityStatus(){
			state := wld_view.States[alive_elevator]
			states[alive_elevator] = HRAElevState{
				Behavior:    state.Behaviour,
				Floor:       state.Floor,
				Direction:   state.Direction,
				CabRequests: state.GetCabRequests(),
			}
		}
	}

	if len(states) == 0{
		return
	}

	input := HRAInput{
		HallRequests: wld_view.GetHallRequests(),
		States:       states,
	}

	// input.PrintInput()
	// wld_view.PrintWorldView()


	jsonBytes, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	ret, err := exec.Command(hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		panic(err)
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		panic(err)
	}

	wld_view.AssignedOrders = *output
	// wld_view.PrintWorldView()
}

func (inp HRAInput) PrintInput() {
	for i, floor := range inp.HallRequests {
		fmt.Printf("Requests at floor: %d\n", i)
		for j, button := range floor {
			fmt.Printf("    Button %d: %t\n", j, button)
		}
	}

	for IP, elev := range inp.States {
		fmt.Printf("State of elevator %s: \n", IP)
		fmt.Printf("    Behaviour: %s\n", elev.Behavior)
		fmt.Printf("    Floor: %d\n", elev.Floor)
		fmt.Printf("    Direction: %s\n", elev.Direction)
		fmt.Printf("    Cab requests: ")
		for _, value := range elev.CabRequests {
			fmt.Printf(" %t", value)
		}
		fmt.Println("")
	}
}
