package process_pair_timer

import (
	"Sanntid/timer"
)


func CheckProcessPairTimeout(tmr *timer.Timer, timer_duration float64, timeout chan<- bool){
	for {
		if tmr.TimerTimedOut(timer_duration){
			timeout<-true
			return
		}
	}
}