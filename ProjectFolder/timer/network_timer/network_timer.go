package network_timer

import (
	"Sanntid/message_handler"
	"Sanntid/timer"
	"Sanntid/world_view"
	"time"
)

func CheckNetworkTimeout(tmr *timer.Timer, worldView *world_view.WorldView, myIP string, msgRx chan <- message_handler.StandardMessage, net_lost chan <- bool) {
	for {
		if tmr.Timer_timedOut(timer.NETWORK_TIMER_TimoutTime) {
			var sendTime string = time.Now().String()[11:19]
			msgRx <-  message_handler.CreateStandardMessage(*worldView, myIP, sendTime)
			net_lost <- true
		}
	}
}