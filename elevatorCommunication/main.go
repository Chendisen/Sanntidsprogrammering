package main

import (
	"Communication/broadcast"
	"Communication/receive"
	"Communication/utilities"
)

func main() {

	IP, _ := broadcast.GetOutboundIP()

	var currentAlive map[string]utilities.Pair

	go broadcast.Peer_broadcastAlive(IP.String())
	go receive.Peer_receiveAlive(&currentAlive)

	for {
		continue
	}
}
