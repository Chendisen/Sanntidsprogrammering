package main

import (
	//"Communication/broadcast"
	"Communication/receive"
)

func main() {
	//go broadcast.Peer_broadcastAlive()
	go receive.Peer_receiveAlive()

	for {
		continue
	}
}
