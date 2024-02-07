package elevator_io

import (
	"Sanntid/driver"
)

const N_FLOORS int = 4
const N_BUTTONS int = 3

const addr string = "localhost"
const port string = "15657"

// type Dirn int

// const (
// 	D_Down Dirn = -1
// 	D_Stop Dirn = 0
// 	D_Up   Dirn = 1
// )

// type Button int

// const (
// 	B_HallUp Button = iota
// 	B_HallDown
// 	B_Cab
// )

type ElevInputDevice struct {
	FloorSensor   func() int
	RequestButton func(int, driver.ButtonType) bool
	StopButton    func() bool
	Obstruction   func() bool
}

type ElevOutputDevice struct {
	FloorIndicator     func(int)
	RequestButtonLight func(int, driver.ButtonType, bool)
	DoorLight          func(bool)
	StopButtonLight    func(bool)
	MotorDirection     func(driver.MotorDirection)
}

func Init() {
	driver.Init(addr+":"+port, N_FLOORS)
}

func wrap_requestButton(f int, b driver.ButtonType) bool {
	return driver.GetButton(b, f)
}

func wrap_requestButtonLight(f int, b driver.ButtonType, v bool) {
	driver.SetButtonLamp(b, f, v)
}

func wrap_motorDirection(d driver.MotorDirection) {
	driver.SetMotorDirection(d)
}

func Elevio_getInputDevice() ElevInputDevice {
	return ElevInputDevice{
		FloorSensor:   driver.GetFloor,
		RequestButton: wrap_requestButton,
		StopButton:    driver.GetStop,
		Obstruction:   driver.GetObstruction,
	}
}

func Elevio_getOutputDevice() ElevOutputDevice {
	return ElevOutputDevice{
		FloorIndicator:     driver.SetFloorIndicator,
		RequestButtonLight: wrap_requestButtonLight,
		DoorLight:          driver.SetDoorOpenLamp,
		StopButtonLight:    driver.SetStopLamp,
		MotorDirection:     wrap_motorDirection,
	}
}

func Elevio_dirn_toString(d driver.MotorDirection) string {
	switch d {
	case driver.MD_Up:
		return "MD_Up"
	case driver.MD_Down:
		return "MD_Down"
	case driver.MD_Stop:
		return "MD_Stop"
	default:
		return "MD_UNDEFINED"
	}
}

func Elevio_button_toString(b driver.ButtonType) string {
	switch b {
	case driver.BT_HallUp:
		return "BT_HallUp"
	case driver.BT_HallDown:
		return "BT_HallDown"
	case driver.BT_Cab:
		return "BT_Cab"
	default:
		return "BT_UNDEFINED"
	}
}
