package elevator

import (
	"fmt"
	"Sanntid/elevator_io"
	"Sanntid/driver"
	// "Sanntid/timer"
)

type ElevatorBehaviour int64 

const (
	EB_Idle 		ElevatorBehaviour = iota
	EB_DoorOpen 	
	EB_Moving		
)

type ClearRequestVariant int64

const (
	// Assume everyone waiting for the elevator gets on the elevator, 
	// they will be traveling in the "wrong" direction for a while
	CV_all		ClearRequestVariant = iota

	// Assume only those that want to travel in the current direction
	// enter the elevator, and keep waiting outside otherwise
	CV_InDirn
)

type Configuration struct {
	ClearRequestVariant 	ClearRequestVariant
	DoorOpenDuration_s 		float64
}

type Elevator struct {
	Floor													int
	Dirn 													driver.MotorDirection
	Request[elevator_io.N_FLOORS][elevator_io.N_BUTTONS]    int
	Behaviour 												ElevatorBehaviour
	Config 													Configuration
}

func eb_toString(eb ElevatorBehaviour) string {
	switch eb {
	case EB_Idle:
		return "EB_Idle"
	case EB_DoorOpen:
		return "EB_DoorOpen"
	case EB_Moving:
		return "EB_Moving"
	default:
		return "EB_UNDEFINED"
	}
}

func Elevator_print(es Elevator) {
	fmt.Println("  +-----------------------+")
	fmt.Printf("  |floor = %2d          |\n  |dirn  = %12s|\n  |behav = %12s|\n", es.Floor, elevator_io.Elevio_dirn_toString(es.Dirn),eb_toString(es.Behaviour))
	fmt.Println("  +-----------------------+")
	fmt.Println("  | up | dn | cab |")
	for floor := elevator_io.N_FLOORS - 1; floor >= 0; floor -- {
		fmt.Printf("  | %d", floor)
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			if ((floor == elevator_io.N_FLOORS && btn == int(driver.BT_HallUp)) ||
				(floor == 0 && btn == int(driver.BT_HallDown))) {
				fmt.Println("|     ")
			} else {
				switch es.Request[floor][btn] {
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
		Floor: 		-1, 
		Dirn: 		driver.MD_Stop,
		Behaviour: 	EB_Idle,
		Config: 	Configuration 	{ClearRequestVariant: 	CV_InDirn,
					 				DoorOpenDuration_s: 	3.0},
	}
}

