package communication

import (
	"Sanntid/communication/bcast"
	"Sanntid/communication/peers"
	"Sanntid/timer"
	"Sanntid/timer/network_timer"
	"Sanntid/world_view"
	"fmt"
	"time"
)

const PeriodInMilliseconds int = 100

func StartCommunication(myView *world_view.WorldView, networkOverview *world_view.NetworkOverview, msg_received chan world_view.StandardMessage, heardFromList *world_view.HeardFromList, ord_updated chan<- bool, wld_updated chan<- bool) {

	time.Sleep(2 * time.Second)

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(55555, networkOverview.GetMyIP(), peerTxEnable)
	go peers.Receiver(55555, peerUpdateCh)

	msgTx := make(chan world_view.StandardMessage, 10)
	msgRx := make(chan world_view.StandardMessage, 10)

	go bcast.Transmitter(11111, msgTx)
	go bcast.Receiver(11111, msgRx)

	var standardMessage world_view.StandardMessage = world_view.CreateStandardMessage(*myView, networkOverview.GetMyIP(), time.Now().String()[11:19])

	var timerNetwork timer.Timer = timer.TimerUninitialized()
	net_lost := make(chan bool)
	go network_timer.CheckNetworkTimeout(&timerNetwork, myView, networkOverview.GetMyIP(), msgRx, net_lost)

	go func() {
		for {
			standardMessage.WorldView = *myView
			standardMessage.SendTime = time.Now().String()[11:19]
			msgTx <- standardMessage
			time.Sleep(time.Duration(PeriodInMilliseconds) * time.Millisecond)
		}
	}()

	go func () {peerUpdateCh <- peers.PeerUpdate{
		Peers: make([]string, 0),
		New: "",
		Lost: make([]string, 0),
	}
	} ()
	
	fmt.Println("Started communications")
	
	for {
		select {
		case p := <-peerUpdateCh:

			if networkOverview.NetworkLost(p) {
				p.Peers = append(p.Peers, networkOverview.GetMyIP())
				timerNetwork.TimerStart(timer.NETWORK_TIMER_TimoutTime)
			} else {
				timerNetwork.TimerStop()
			}

			networkOverview.UpdateNetworkOverview(p)
			if len(p.New) > 0 {
				heardFromList.AddNodeToList(p.New)
			}
			if len(p.Lost) > 0 {
				wld_updated <- true
			}

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			fmt.Printf((" Am i master?:  %t\n"), networkOverview.AmIMaster())

		case recievedMsg := <-msgRx:
			//myView.UpdateWorldView(recievedMsg.WorldView, recievedMsg.IPAddress, recievedMsg.SendTime, networkOverview.MyIP, *networkOverview, heardFromList, lightArray, ord_updated, wld_updated)
			msg_received <- recievedMsg
		case networkLost := <-net_lost:
			if networkLost {
				timerNetwork.TimerStart(timer.NETWORK_TIMER_TimoutTime)
			} else {
				timerNetwork.TimerStop()
			}
		}
	}
}
