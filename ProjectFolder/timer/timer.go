package timer

import (
	"math"
	"time"
)

type RequestType int 
type TimerType int

const (
	Start RequestType = iota
	Stop
	TimedOut
)

const (
	DoorTimer TimerType = iota
	NetworkTimer
	ProcessPairTimer
	WatchdogTimer
)

type TimerRequest struct{
	RequestType RequestType
	TimerType TimerType
}

func GenerateTimerRequest(reqType RequestType, tmrType TimerType) TimerRequest{
	return TimerRequest{RequestType: reqType, TimerType: tmrType}
}

const DOOR_OPEN_TimeoutTime float64 = 3
const WATCHDOG_TimeoutTime float64 = 5
const PROCESS_PAIR_TimeoutTime float64 = 3
const NETWORK_TIMER_TimoutTime float64 = 0.5

type Timer struct {
	TimerEndTime float64
	timerActive  bool
}

func Timer_uninitialized() Timer {
	return Timer{TimerEndTime: 0, timerActive: false}
}

func Get_current_time() float64 {
	return (float64(time.Now().Second()) + float64(time.Now().Nanosecond())*float64(0.000000001))
}

func (tmr *Timer) Timer_start(duration float64) {
	tmr.TimerEndTime = math.Mod((Get_current_time() + duration), 60.0)
	tmr.timerActive = true
}

func (tmr *Timer) Timer_stop() {
	tmr.timerActive = false
}

func (tmr *Timer) Timer_timedOut(timer_duration float64) bool {
	return (tmr.timerActive && (Get_current_time() > tmr.TimerEndTime) && !(tmr.TimerEndTime < timer_duration && Get_current_time() > (60 - timer_duration)))
}