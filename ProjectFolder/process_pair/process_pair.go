package process_pair

import (
	//"Sanntid/message_handler"
	"Sanntid/communication/bcast"
	"Sanntid/communication/peers"
	"Sanntid/timer"
	"Sanntid/timer/process_pair_timer"
	"Sanntid/world_view"
	"fmt"
	"time"
)

func ProcessPair(myIP string, storedView *world_view.WorldView, tmr *timer.Timer, startNew chan<- bool) {

	time.Sleep(2 * time.Second)

	peerUpdateCh := make(chan peers.PeerUpdate)
	go peers.Receiver(55555, peerUpdateCh)

	msgRx := make(chan world_view.StandardMessage, 10)
	go bcast.Receiver(11111, msgRx)

	var p peers.PeerUpdate

	fmt.Println("Started listening to primary")

	timeOut := make(chan bool)
	tmr.Timer_start(timer.PROCESS_PAIR_TimeoutTime)
	go process_pair_timer.CheckProcessPairTimeout(tmr, timer.PROCESS_PAIR_TimeoutTime, timeOut)

	for {
		select {
		case p = <-peerUpdateCh:
			fmt.Println("Peer updated")

			if len(p.Lost) > 0  {
				for _,IP := range p.Lost {
					if IP == myIP {
						tmr.Timer_start(timer.PROCESS_PAIR_TimeoutTime)
						break
					}
				}
				} else if p.New == myIP {
					tmr.Timer_stop()
				}
				
			fmt.Printf("new peer: %s My IP: %s", p.New, myIP)
				
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case recievedMsg := <-msgRx:
			if recievedMsg.IPAddress == myIP {
				*storedView = recievedMsg.WorldView
			}

		case <-timeOut:
			fmt.Println("Are we often in here?")
			if len(p.Peers) > 0{
				*storedView = world_view.MakeWorldView(myIP)
			}
			startNew<-true
			return
		}
	}
}

