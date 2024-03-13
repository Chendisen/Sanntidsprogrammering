package timer

import (
	"math"
	"time"
)

type Timer struct {
	timerEndTime float64
	timerActive  bool
}

func Timer_uninitialized() Timer {
	return Timer{timerEndTime: 0, timerActive: false}
}

func get_current_time() float64 {
	return (float64(time.Now().Second()) + float64(time.Now().Nanosecond())*float64(0.000000001))
}

func Timer_start(tmr *Timer, duration float64) {
	tmr.timerEndTime = math.Mod((get_current_time() + duration), 60.0)
	tmr.timerActive = true
}

func Timer_stop(tmr *Timer) {
	tmr.timerActive = false
}

func Timer_timedOut(tmr *Timer) bool {
	return (tmr.timerActive && (get_current_time() > tmr.timerEndTime) && !(tmr.timerEndTime < 3.0 && get_current_time() > 57))
}

func (tmr *Timer) TimeOut(timeOut chan<- bool){
	for {
		if Timer_timedOut(tmr){
			timeOut<-true
			return
		}
	}
}