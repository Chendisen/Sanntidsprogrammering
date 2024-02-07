package timer

import (
	"time"
)

func get_current_time() float64 {
	return (time.Now().Seconds() + time.Now().Milliseconds()*0.000001)
}

var timerEndTime float64
var timerActive bool

func timer_start(duration float64) {
	timerEndTime = get_current_time() + duration
	timerActive = 1
}

func timer_stop() {
	timerActive = 0
}

func timer_timedOut() bool {
	return (timerActive && (get_current_time() > timerEndTime))
}
