package world_view

import(
	"Sanntid/driver"
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

func (lightArray *LightArray) InitLights(myIP string, worldView WorldView){
	for floor, buttons := range worldView.GetHallRequests() {
		for button, value := range buttons {
			(*lightArray)[floor][button] = value
		}
	}
	for floor,value := range worldView.GetMyCabRequests(myIP) {
		(*lightArray)[floor][driver.BT_Cab] = value
	}
}
