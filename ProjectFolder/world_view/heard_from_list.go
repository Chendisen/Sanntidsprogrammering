package world_view

import (
	"Sanntid/resources/driver"
	"fmt"
)


type HeardFromList struct {
	HeardFrom map[string][][3]bool
}

func MakeHeardFromList(myIP string) HeardFromList {
	heardFromList := HeardFromList{HeardFrom: make(map[string][][3]bool)}
	heardFromList.HeardFrom[myIP] = make([][3]bool, driver.N_FLOORS)
	return heardFromList
}

func (heardFromList HeardFromList) ShouldResetAtFloorButton(floor int, button int, networkOverview NetworkOverview) bool {
	var count int = 0
	for _, buttonArray := range heardFromList.HeardFrom {
		if buttonArray[floor][button] {
			count++
		}
	}
	return count == len(networkOverview.NodesAlive)
}

func (heardFromList HeardFromList) ShouldAddNode(ip string) bool {
	var check bool = true
	for IP := range heardFromList.HeardFrom {
		if IP == ip {
			check = false
			return check
		}
	}
	return check
}

func (heardFromList *HeardFromList) SetHeardFrom(networkOverview NetworkOverview, msgIP string, floor int, button int) {
	for _, id := range networkOverview.NodesAlive {
		if id == msgIP {
			heardFromList.HeardFrom[msgIP][floor][button] = true
			return
		}
	}
}

func (heardFromList *HeardFromList) GetHeardFrom(networkOverview NetworkOverview, msgIP string, floor int, button int) bool {
	for _, id := range networkOverview.NodesAlive {
		if id == msgIP {
			return heardFromList.HeardFrom[msgIP][floor][button]
		}
	}
	return false
}

func (heardFromList *HeardFromList) CheckHeardFromAll(networkOverview NetworkOverview, floor int, button int) bool {
	var heard_from_all bool = true
	for _, alv_nodes := range networkOverview.NodesAlive {
		heard_from_all = heard_from_all && heardFromList.HeardFrom[alv_nodes][floor][button]
	}
	return heard_from_all
}

func (heardFromList *HeardFromList) ClearHeardFrom(floor int, button int) {
	for _, hfl_buttons := range heardFromList.HeardFrom {
		hfl_buttons[floor][button] = false
	}
}

func (heardFromList *HeardFromList) AddNodeToList(newIP string) {
	heardFromList.HeardFrom[newIP] = make([][3]bool, driver.N_FLOORS)
}

func (heardFromList HeardFromList) Print() {
	fmt.Println("We have heard from: ")
	for IP := range heardFromList.HeardFrom {
		fmt.Printf("	%s\n", IP)
	}
	fmt.Printf("")

	for IP, table := range heardFromList.HeardFrom {
		fmt.Printf("Elevator: %s \n", IP)
		for floor, buttons := range table {
			fmt.Printf("	Floor: %d", floor)
			for button := range buttons {
				fmt.Printf("	Button: %d", button)
			}
			fmt.Print("\n")
		}
	}
}
