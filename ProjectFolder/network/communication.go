package network

import (
	//"Sanntid/network/bcast"
	"Sanntid/network/peers"
	"Sanntid/world_view"
	"fmt"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//
//	will be received as zero-values.

func StartCommunication(myIP string, c chan world_view.WorldView, al *world_view.AliveList){ //, myView *world_view.WorldView, ord_updated chan<- int) {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, myIP, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	//helloTx := make(chan world_view.WorldView)
	//helloRx := make(chan world_view.WorldView)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	//go bcast.Transmitter(16569, helloTx)
	//go bcast.Receiver(16569, helloRx)

	// The example message. We just send one of these every second.

	/*go func() {
		var wv world_view.WorldView
		for {
			wv = <-c
			helloTx <- wv
		}
	}()*/

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:

			al.UpdateAliveList(p)

			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		/*case newView := <-helloRx:
			myView.UpdateWorldView(newView, "123", al.MyIP, *al, ord_updated)
		*/}
	}
}
