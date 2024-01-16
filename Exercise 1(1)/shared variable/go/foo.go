// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
)

var i int = 0
var q int = 0

func incrementing(c, quit chan int) {
	//TODO: increment i 1000000 times
	for j := 0; j < 1000000; j++ {
		c <- 1
	}
	quit <- 1
}

func decrementing(c, quit chan int) {
	//TODO: decrement i 1000000 times
	for j := 0; j < 500000; j++ {
		c <- 1
	}
	quit <- 1
}

func sync(c1, c2, quit chan int) {
	for {
		select {
		case <-c1:
			i += 1
		case <-c2:
			i -= 1
		case <-quit:
			q += 1
			if q == 2 {
				Println(i)
				return
			}
		}
	}
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(2)

	// TODO: Spawn both functions as goroutines
	c1 := make(chan int)
	c2 := make(chan int)
	quit := make(chan int)

	go incrementing(c1, quit)
	go decrementing(c2, quit)
	sync(c1, c2, quit)

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	//time.Sleep(500 * time.Millisecond)
	//Println("The magic number is:", i)

}
