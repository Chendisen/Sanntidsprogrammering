package main

import "Driver-go/elevio"
import "fmt"


func main(){

    numFloors := 4

    elevio.Init("localhost:15657", numFloors)
    
    var d elevio.MotorDirection = elevio.MD_Up
    //elevio.SetMotorDirection(d)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)
    
    
    for {
        select {
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            
        case a := <- drv_floors:
            fmt.Printf("%+v\n", a)
            if a == numFloors-1 {
                d = elevio.MD_Down
            } else if a == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)
            
            
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }
            
        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        }
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