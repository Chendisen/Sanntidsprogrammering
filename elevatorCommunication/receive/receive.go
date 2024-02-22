package receive

import (
	"fmt"
	"net"
)

const port string = "12345"
const messageSize int = 1024

func Peer_receiveAlive() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Println("Failed to resolve UDP address")
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to listen to UDP")
		return
	}
	defer conn.Close()

	buffer := make([]byte, messageSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read UDP message")
			continue
		}

		fmt.Printf("Recieved message: %s", string(buffer[:n]))
	}
}
