package main

import (
	//"Sanntid/driver"
	//"Sanntid/elevator"
	//"Sanntid/fsm"
	//"Sanntid/timer"
	"Sanntid/world_view"
	//"Sanntid/network/localip"
	"Sanntid/network"
	//"fmt"
)

func main() {
	const numFloors int = 4
	//driver.Init("localhost:15657", numFloors)

	//var elev elevator.Elevator = elevator.Elevator_uninitialized()
	//var tmr timer.Timer = timer.Timer_uninitialized()
	var alv_list world_view.AliveList = world_view.MakeAliveList()
	//var wld_view world_view.WorldView = world_view.MakeWorldView(alv_list.MyIP)

	wv_update := make(chan world_view.WorldView)


	go network.StartCommunication(alv_list.MyIP, wv_update, &alv_list)

	for {
		continue
	}

}	