package network

import (
	"Sanntid/message_handler"
	"Sanntid/network/bcast"
	"Sanntid/network/peers"
	"Sanntid/world_view"
	"fmt"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//
//	will be received as zero-values.

func StartCommunication(myIP string, myView *world_view.WorldView, al *world_view.AliveList, hfl *world_view.HeardFromList, lgt_array *[][3]bool, ord_updated chan<- bool, wld_updated chan<- bool) {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network

	time.Sleep(5 * time.Second)

	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(55555, myIP, peerTxEnable)
	go peers.Receiver(55555, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	msgTx := make(chan message_handler.StandardMessage, 10)
	msgRx := make(chan message_handler.StandardMessage, 10)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(11111, msgTx)
	go bcast.Receiver(11111, msgRx)

	// The example message. We just send one of these every second.

	var sm message_handler.StandardMessage = message_handler.CreateStandardMessage(*myView, myIP, time.Now().String()[11:19])

	go func() {
		for {
			sm.WorldView = *myView
			sm.SendTime = time.Now().String()[11:19]
			msgTx <- sm
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:

			al.UpdateAliveList(p)
			if myView.ShouldAddNode(p.New){
				myView.AddNodeToWorldView(p.New)
			}
			myView.LastHeard[p.New] = time.Now().String()[11:19]

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			fmt.Printf((" Am i master?:  %t\n"), (*al).AmIMaster())

		case recievedMsg := <-msgRx:
			myView.UpdateWorldView(recievedMsg.WorldView, recievedMsg.IPAddress, recievedMsg.SendTime, al.MyIP, *al, hfl, lgt_array, ord_updated, wld_updated)
		}
	}
}
