package run_elevator
/*
import (
	"Sanntid/driver"
	"Sanntid/elevator_io"
	"Sanntid/fsm"
	"Sanntid/timer"
	"fmt"
	"time"
)

func Run_elevator() {
	fmt.Printf("Started!\n")

	const inputPollRate_ms int = 25

	var input elevator_io.ElevInputDevice = elevator_io.Elevio_getInputDevice()

	if input.FloorSensor() == -1 {
		fsm.Fsm_onInitBetweenFloors()
	}

	for {
		{ // Request button
			var prev [elevator_io.N_FLOORS][elevator_io.N_BUTTONS]int
			for floor := 0; floor < elevator_io.N_FLOORS; floor++ {
				for btn := 0; btn < elevator_io.N_BUTTONS; btn++ {
					var v bool = input.RequestButton(floor, driver.ButtonType(btn))
					if v && v != driver.IntToBool(prev[floor][btn]) {
						fsm.Fsm_onRequestButtonPress(floor, driver.ButtonType(btn))
					}
					prev[floor][btn] = driver.BoolToInt(v)
				}
			}
		}

		{ // Floor sensor
			var prev int = -1
			var floor int = input.FloorSensor()
			if floor != -1 && floor != prev {
				fsm.Fsm_onFloorArrival(floor)
			}
			prev = floor
		}

		{ // Timer
			if timer.Timer_timedOut() {
				timer.Timer_stop()
				fsm.Fsm_onDoorTimeout()
			}
		}

		time.Sleep(time.Duration(inputPollRate_ms * 1000))
	}
}
*/