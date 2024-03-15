package elevator

import(
	"Sanntid/resources/driver"
)

type LightArray [][3]bool

func MakeLightArray() LightArray {
	return make([][3]bool, driver.N_FLOORS)
}

func (lightArray LightArray) SetAllLights() {
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		for btn := 0; btn < driver.N_BUTTONS; btn++ {
			driver.SetButtonLamp(driver.ButtonType(btn), floor, lightArray[floor][btn])
		}
	}
}

func (lightArray *LightArray) InitLights(hallRequests [][2]bool, cabRequests []bool){
	for floor, buttons := range hallRequests {
		for button, value := range buttons {
			(*lightArray)[floor][button] = value
		}
	}
	for floor,value := range cabRequests {
		(*lightArray)[floor][driver.BT_Cab] = value
	}
}

func (lightArray *LightArray) SetElevatorLight(floor int, button int) {
	(*lightArray)[floor][button] = true
}

func (lightArray *LightArray) ClearElevatorLight(floor int, button int) {
	(*lightArray)[floor][button] = false
}