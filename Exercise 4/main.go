package main

import (
	"fmt"
	"net"
	"os/exec"
	"time"
)

const duration float64 = 1.0

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
	tmr.timerEndTime = get_current_time() + duration
	tmr.timerActive = true
}

func Timer_stop(tmr *Timer) {
	tmr.timerActive = false
}

func Timer_timedOut(tmr *Timer) bool {
	return (tmr.timerActive && (get_current_time() > tmr.timerEndTime))
}

func check_timeout(tmr *Timer) {

	// for {
	// 	receiveUDP()
	// }

	var it int = 0
	Timer_start(tmr, duration)

	for {
		if Timer_timedOut(tmr) {
			it += 1
			fmt.Print(it)
			Timer_start(tmr, duration)
			broadcastUDP()

		}
		if it == 5 {
			it = 0
		}
	}

	exec.Command("gnome-terminal", "--", "go", "run", "").Run()
}

func broadcastUDP() {
	connection1, err := net.Dial("udp", "localhost:12345")
	if err != nil {
		panic(err)
	}
	
	_, err = connection1.Write([]byte("Don't @ me bro!\000"))
	if err != nil {
		panic(err)
	}
	
	connection1.Close()
}

func receiveUDP() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":20014")
	if err != nil {
		panic(err)
	}

	connection, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
}

func main() {
	var tmr Timer = Timer{timerEndTime: 0, timerActive: false}
	check_timeout(&tmr)
}
