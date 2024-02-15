package run_elevator
/*
import (
	"Sanntid/driver"
	"Sanntid/fsm"
    "Sanntid/elevator_io"
	"fmt"
)


func Run_elevator_go(){

    numFloors := 4

    driver.Init("localhost:15657", numFloors)

    var input elevator_io.ElevInputDevice = elevator_io.Elevio_getInputDevice()

	if input.FloorSensor() == -1 {
		fsm.Fsm_onInitBetweenFloors()
	}
    
    //var d driver.MotorDirection = driver.MD_Up
    //driver.SetMotorDirection(d)
    
    drv_buttons := make(chan driver.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)    
    
    go driver.PollButtons(drv_buttons)
    go driver.PollFloorSensor(drv_floors)
    go driver.PollObstructionSwitch(drv_obstr)
    go driver.PollStopButton(drv_stop)
    
    
    for {
        select {
        case a := <- drv_buttons:
            fmt.Printf("%+v\n", a)
            //driver.SetButtonLamp(a.Button, a.Floor, true)
            fsm.Fsm_onRequestButtonPress(a.Floor, a.Button)
            
        case a := <- drv_floors:
            fmt.Printf("%+v\n", a)
            // if a == numFloors-1 {
            //     d = driver.MD_Down
            // } else if a == 0 {
            //     d = driver.MD_Up
            // }
            // driver.SetMotorDirection(d)
            fsm.Fsm_onFloorArrival(a)
            
            
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            // if a {
            //     driver.SetMotorDirection(driver.MD_Stop)
            // } else {
            //     driver.SetMotorDirection(d)
            // }
            fsm.Fsm_onDoorTimeout()
            
        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := driver.ButtonType(0); b < 3; b++ {
                    driver.SetButtonLamp(b, f, false)
                }
            }
        }
    }    
}
*/