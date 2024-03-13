package process_pair

import (
	"Sanntid/message_handler"
	"Sanntid/network/bcast"
	"Sanntid/network/peers"
	"Sanntid/timer"
	"Sanntid/world_view"
	"fmt"
	"time"
)


func ProcessPair(myIP string, storedView *world_view.WorldView, tmr *timer.Timer, startNew chan<- bool) {

	time.Sleep(2 * time.Second)

	peerUpdateCh := make(chan peers.PeerUpdate)
	go peers.Receiver(55555, peerUpdateCh)

	msgRx := make(chan message_handler.StandardMessage, 10)
	go bcast.Receiver(11111, msgRx)

	var p peers.PeerUpdate

	fmt.Println("Started communications")

	timeOut := make(chan bool)
	timer.Timer_start(tmr, 3)
	go tmr.TimeOut(timeOut)

	for {
		select {
		case p = <-peerUpdateCh:

			if len(p.Lost) > 0  {
				for _,IP := range p.Lost {
					if IP == myIP {
						timer.Timer_start(tmr, 3)
						break
					}
				}
			} else if p.New == myIP {
				timer.Timer_stop(tmr)
			}
			

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case recievedMsg := <-msgRx:
			if recievedMsg.IPAddress == myIP {
				*storedView = recievedMsg.WorldView
			}

		case <-timeOut:
			if len(p.Peers) > 0{
				*storedView = world_view.MakeWorldView(myIP)
			}
			startNew<-true
			return
		}
	}
}

