package elevator

import (
	"Sanntid/driver"
	"fmt"
)

type ElevatorBehaviour int64 

const (
	EB_Idle 		ElevatorBehaviour = iota
	EB_DoorOpen 	
	EB_Moving		
)

type Configuration struct {
	DoorOpenDuration_s 		float64
}

type Elevator struct {
	Floor													int
	Dirn 													driver.MotorDirection
	Request[driver.N_FLOORS][driver.N_BUTTONS]    			int
	Behaviour 												ElevatorBehaviour
	Config 													Configuration
	DoorObstructed											bool
}

func ElevatorBehaviourToString(elevatorBehaviour ElevatorBehaviour) string {
	switch elevatorBehaviour {
	case EB_Idle:
		return "idle"
	case EB_DoorOpen:
		return "doorOpen"
	case EB_Moving:
		return "moving"
	default:
		return "undefined"
	}
}

func Elevator_print(elev Elevator) {
	fmt.Println("  +-----------------------+")
	fmt.Printf("  |floor = %2d          |\n  |dirn  = %12s|\n  |behav = %12s|\n", elev.Floor, driver.DriverDirectionToString(elev.Dirn), ElevatorBehaviourToString(elev.Behaviour))
	fmt.Println("  +-----------------------+")
	fmt.Println("  | up | dn | cab |")
	for floor := driver.N_FLOORS - 1; floor >= 0; floor -- {
		fmt.Printf("  | %d", floor)
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			if ((floor == driver.N_FLOORS && btn == int(driver.BT_HallUp)) ||
				(floor == 0 && btn == int(driver.BT_HallDown))) {
				fmt.Println("|     ")
			} else {
				switch elev.GetElevatorRequest(floor, btn) {
				case 1:
					fmt.Println("|  #  ")
				case 0:
					fmt.Println("|  -  ")
				}
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  +-----------------------+")
}

func Elevator_uninitialized() Elevator {
	return Elevator {
		Floor: 			-1, 
		Dirn: 			driver.MD_Stop,
		Behaviour: 		EB_Idle,
		Config: 		Configuration 	{DoorOpenDuration_s: 	3.0},
		DoorObstructed: false,
	}
}

func (elev *Elevator) GetElevatorRequest(floor int, button int) int {
	return (*elev).Request[floor][button]
}

func (elev *Elevator) SetElevatorRequest(floor int, button int, value int) {
	(*elev).Request[floor][button] = value
}