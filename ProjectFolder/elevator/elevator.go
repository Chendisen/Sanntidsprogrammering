package elevator

import (
	"fmt"
	"Sanntid/elevator_io"
	"Sanntid/timer"
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

type Config struct {
	clearRequestVariant 	ClearRequestVariant
	doorOpenDuration_s 		float64
}

type Elevator struct {
	floor													int64
	dirn 													elevator_io.Dirn
	request[elevator_io.N_FLOORS][elevator_io.N_BUTTONS] 	int
	behaviour 												ElevatorBehaviour
	config 													Config
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
	fmt.Println("  +-----------------------+\n")
	fmt.Println(
		"  |floor = %2d          |\n"
		"  |dirn  = %12s|\n"
		"  |behav = %12s|\n",
		es.floor, 
		elevator_io.elevio_dirn_toString(es.dirn),
		eb_toString(es.behaviour)
	)
	fmt.Println("  +-----------------------+\n")
	fmt.Println("  | up | dn | cab |\n")
	for floor := elevator_io.N_FLOORS - 1; floor >= 0; floor -- {
		fmt.Println("  | %d", floor)
		for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
			if ((floor == elevator_io.N_FLOORS && btn == elevator_io.B_HallUp) ||
				(floor == 0 && btn == elevator_io.B_Halldown)
			) {
				fmt.Println("|     ")
			} else {
				switch es.request[floor][btn] {
				case 1:
					fmt.Println("|  #  ")
				case 0:
					fmt.Println("|  -  ")
				}
			}
		}
		fmt.Println("|\n")
	}
	fmt.Println("  +-----------------------+\n")
}

func Elevator_uninitialized() Elevator {
	return Elevator {
		floor: 		-1, 
		dirn: 		elevator_io.D_Stop,
		behaviour: 	EB_Idle
		config: 	{clearRequestVariant: 	CV_InDirn
					 doorOpenDuration_s: 	3.0},
	}
}