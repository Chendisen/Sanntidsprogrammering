package broadcast

import (
	"fmt"
	"net"
	"time"
)

const broadcastAddr string = "255.255.255.255"
const port string = "12345"

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func Peer_broadcastAlive() {
	ip, err := GetOutboundIP()
	if err != nil {
		fmt.Println("Failed to get own ip")
		return
	}

	fmt.Print(ip)

	udpAddr, err := net.ResolveUDPAddr("udp", broadcastAddr+":"+port)
	if err != nil {
		fmt.Println("Failed to resolve udp address")
		return
	}

	for{
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			fmt.Println("Failed to dial udp")
		}

		message := []byte(ip)

		_, err = conn.Write(message)
		if err != nil {
			fmt.Println("Error sending UDP message: ", err)
		}

		fmt.Println("UDP message sent successfully")

		conn.Close()
		time.Sleep(time.Millisecond * 100)
	}
}


