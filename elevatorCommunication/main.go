package main

import (
	"Communication/broadcast"
	"Communication/recieve"
)

func main() {
	go broadcast.Peer_broadcastAlive()
	go recieve.Peer_recieveAlive()
}
