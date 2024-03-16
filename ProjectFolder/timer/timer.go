package timer

import (
	"math"
	"time"
)

const (
	DoorTimer TimerType = iota
	NetworkTimer
	ProcessPairTimer
	WatchdogTimer
)

const DOOR_OPEN_TimeoutTime float64 = 3
const WATCHDOG_TimeoutTime float64 = 5
const PROCESS_PAIR_TimeoutTime float64 = 3
const NETWORK_TIMER_TimoutTime float64 = 0.5

type Timer struct {
	timerEndTime float64
	timerActive  bool
}

type RequestType int 
type TimerType int

const (
	Start RequestType = iota
	Stop
	TimedOut
)

type TimerRequest struct{
	RequestType RequestType
	TimerType TimerType
}


func GenerateTimerRequest(reqType RequestType, tmrType TimerType) TimerRequest{
	return TimerRequest{RequestType: reqType, TimerType: tmrType}
}


func TimerUninitialized() Timer {
	return Timer{timerEndTime: 0, timerActive: false}
}

func getCurrentTime() float64 {
	return (float64(time.Now().Second()) + float64(time.Now().Nanosecond())*float64(0.000000001))
}

func (tmr *Timer) TimerStart(duration float64) {
	tmr.timerEndTime = math.Mod((getCurrentTime() + duration), 60.0)
	tmr.timerActive = true
}

func (tmr *Timer) TimerStop() {
	tmr.timerActive = false
}

func (tmr *Timer) TimerTimedOut(timer_duration float64) bool {
	return (tmr.timerActive && (getCurrentTime() > tmr.timerEndTime) && !(tmr.timerEndTime < timer_duration && getCurrentTime() > (60 - timer_duration)))
}