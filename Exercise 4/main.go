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

}

func broadcastUDP() {
	for {
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
}

func createBackup() {
	exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
}

func receiveUDP(c chan string, quit chan int) {

	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:12345")
	if err != nil {
		panic(err)
	}

	connection, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1024)

	var i int = 0
	for {
		i += 1
		select {
		case <-quit:
			connection.Close()
			fmt.Print("Help me")
			return
		default:
			mLen, err := connection.Read(buffer)
			if err != nil {
				fmt.Println("Error reading: ", err.Error())
			}
			fmt.Println("Received: ", string(buffer[:mLen]), i)

			c <- string(buffer[:mLen])
		}
	}
}

func main() {
	var isMain bool = false

	var recievedMessage string = ""
	var recieveMessageTimer Timer = Timer{timerEndTime: 0, timerActive: false}

	c := make(chan string, 5)
	quit := make(chan int, 5)
	go receiveUDP(c, quit)

	Timer_start(&recieveMessageTimer, 3)

	for !isMain {

		fmt.Print(Timer_timedOut((&recieveMessageTimer)), recievedMessage)

		select {
		case <-c:
			recievedMessage = <-c
		default:
			if recievedMessage == "" && Timer_timedOut(&recieveMessageTimer) {
				Timer_stop(&recieveMessageTimer)
				isMain = true
				quit <- 1
			}
			if recievedMessage != "" {
				Timer_start(&recieveMessageTimer, 3)
			}
			recievedMessage = ""
		}
	}

	time.Sleep(time.Millisecond * 1000)
	fmt.Println("Creating backup!")
	createBackup()

	go broadcastUDP()
}
