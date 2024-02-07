package elevator_io

import (
	"fmt"
	""
)

const N_FLOORS int = 4
const N_BUTTONS int = 3

type Dirn int64

const (
	D_Down Dirn = -1
	D_Stop Dirn = 0
	D_Up   Dirn = 1
)

type Button int64

const (
	B_HallUp Button = iota
	B_Halldown
	B_Cab
)

type ElevInputDevice struct {
	floorSensor   func() int
	requestButton func(int, Button) int
	stopButton    func() int
	obstruction   func() int
}

type ElevOutputDevice struct {
	floorIndicator     func(int)
	requestButtonLight func(int, Button, int)
	doorLight func(int)
	stopButtonLight func(int)
	motorDirection func(Dirn)
}

func elevio_getInputDevice() ElevInputDevice {

}

func elevio_dirn_toString(d Dirn) string {

}

func elevio_button_toString(b Button) string
{

}
