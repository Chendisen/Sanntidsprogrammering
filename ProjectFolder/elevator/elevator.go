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

