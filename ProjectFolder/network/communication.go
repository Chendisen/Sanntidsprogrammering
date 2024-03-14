package network

import (
	"Sanntid/message_handler"
	"Sanntid/network/bcast"
	"Sanntid/network/peers"
	"Sanntid/world_view"
	"fmt"
	"time"
)

func StartCommunication(myIP string, myView *world_view.WorldView, al *world_view.NetworkOverview, hfl *world_view.HeardFromList, lightArray *[][3]bool, ord_updated chan<- bool, wld_updated chan<- bool) {

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

	go func() {
		for {
			sm.WorldView = *myView
			sm.SendTime = time.Now().String()[11:19]
			msgTx <- sm
			time.Sleep(100 * time.Millisecond)
		}
	}()

	fmt.Println("Started communications")
	for {
		select {
		case p := <-peerUpdateCh:

			al.UpdateNetworkOverview(p)
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
			fmt.Printf((" Am i master?:  %t\n"), (*al).AmIMaster())

		case recievedMsg := <-msgRx:
			myView.UpdateWorldView(recievedMsg.WorldView, recievedMsg.IPAddress, recievedMsg.SendTime, al.MyIP, *al, hfl, lightArray, ord_updated, wld_updated)
		}
	}
}
