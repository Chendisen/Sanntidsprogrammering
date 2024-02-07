package timer

import (
	"time"
)

func get_current_time() float64 {
	return (float64(time.Now().Second()) + float64(time.Now().Nanosecond())*float64(0.000000001))
}

var timerEndTime float64
var timerActive bool

func Timer_start(duration float64) {
	timerEndTime = get_current_time() + duration
	timerActive = true
}

func Timer_stop() {
	timerActive = false
}

func Timer_timedOut() bool {
	return (timerActive && (get_current_time() > timerEndTime))
}
