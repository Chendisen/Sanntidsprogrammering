package network_timer

import (
	"Sanntid/timer"
	"Sanntid/world_view"
	"time"
)

func CheckNetworkTimeout(tmr *timer.Timer, worldView *world_view.WorldView, myIP string, msgRx chan <- world_view.StandardMessage, net_lost chan <- bool) {
	for {
		if tmr.TimerTimedOut(timer.NETWORK_TIMER_TimoutTime) {
			var sendTime string = time.Now().String()[11:19]
			msgRx <-  world_view.CreateStandardMessage(*worldView, myIP, sendTime)
			net_lost <- true
		}
	}
}