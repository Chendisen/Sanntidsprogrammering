package network

import (
	"Sanntid/message_handler"
	"Sanntid/network/bcast"
	"Sanntid/network/peers"
	"Sanntid/world_view"
	"Sanntid/timer"
	"Sanntid/network_timer"
	"fmt"
	"time"
)

func StartCommunication(myIP string, myView *world_view.WorldView, networkOverview *world_view.NetworkOverview, hfl *world_view.HeardFromList, lightArray *[][3]bool, ord_updated chan<- bool, wld_updated chan<- bool) {

	time.Sleep(2 * time.Second)

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(55555, myIP, peerTxEnable)
	go peers.Receiver(55555, peerUpdateCh)

	msgTx := make(chan message_handler.StandardMessage, 10)
	msgRx := make(chan message_handler.StandardMessage, 10)

	go bcast.Transmitter(11111, msgTx)
	go bcast.Receiver(11111, msgRx)

	var sm message_handler.StandardMessage = message_handler.CreateStandardMessage(*myView, myIP, time.Now().String()[11:19])

	var timerNetwork timer.Timer = timer.Timer_uninitialized()
	net_lost := make(chan bool)
	go network_timer.NetworkTimer(&timerNetwork, myView, networkOverview.MyIP, msgRx, net_lost)


	go func() {
		for {
			sm.WorldView = *myView
			sm.SendTime = time.Now().String()[11:19]
			msgTx <- sm
			time.Sleep(100 * time.Millisecond)
		}
	}()

	fmt.Println("Started communications")
	// go func () {peerUpdateCh <- peers.PeerUpdate{
	// 		Peers: make([]),
	// 		New: "",
	// 		Lost: make([]string, 1),
	// 	} 
	// } ()
	// fmt.Println("Are we past here?")

	for {
		select {
		case p := <-peerUpdateCh:

			if networkOverview.NetworkLost(p) {
				p.Peers = append(p.Peers, networkOverview.MyIP)
				timerNetwork.Timer_start(timer.NETWORK_TIMER_TimoutTime)
			} else {
				timerNetwork.Timer_stop()
			}

			networkOverview.UpdateNetworkOverview(p)
			if len(p.New) > 0 {
				if myView.ShouldAddNode(p.New){
					myView.AddNodeToWorldView(p.New)
				}
				hfl.AddNodeToList(p.New)
			}

			if len(p.Lost) > 0 {
				wld_updated<-true
			}

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			fmt.Printf((" Am i master?:  %t\n"), networkOverview.AmIMaster())

			// // Safe to assume that you are always online even if network fails
			// if !networkOverview.AmIAlive(networkOverview.GetMyIP()) {
			// 	myView.AddNodeToWorldView(networkOverview.GetMyIP())
			// }

		case recievedMsg := <-msgRx:
			myView.UpdateWorldView(recievedMsg.WorldView, recievedMsg.IPAddress, recievedMsg.SendTime, networkOverview.MyIP, *networkOverview, hfl, lightArray, ord_updated, wld_updated)
		case networkLost := <- net_lost:
			if (networkLost) {
				timerNetwork.Timer_start(timer.NETWORK_TIMER_TimoutTime)
			} else {
				timerNetwork.Timer_stop()
			}
		}
	}
}
