package timer

import (
	"math"
	"time"
)

const DOOR_OPEN_TimeoutTime float64 = 3
const WATCHDOG_TimeoutTime float64 = 5
const PROCESS_PAIR_TimeoutTime float64 = 3

type Timer struct {
	timerEndTime float64
	timerActive  bool
}

func Timer_uninitialized() Timer {
	return Timer{timerEndTime: 0, timerActive: false}
}

func Get_current_time() float64 {
	return (float64(time.Now().Second()) + float64(time.Now().Nanosecond())*float64(0.000000001))
}

func (tmr *Timer) Timer_start(duration float64) {
	tmr.timerEndTime = math.Mod((Get_current_time() + duration), 60.0)
	tmr.timerActive = true
}

func (tmr *Timer) Timer_stop() {
	tmr.timerActive = false
}

func (tmr *Timer) Timer_timedOut(timer_duration float64) bool {
	return (tmr.timerActive && (Get_current_time() > tmr.timerEndTime) && !(tmr.timerEndTime < timer_duration && Get_current_time() > (60 - timer_duration)))
}

func (tmr *Timer) TimeOut(timer_duration float64, timeout chan<- bool){
	for {
		if tmr.Timer_timedOut(timer_duration){
			timeout<-true
			return
		}
	}
}