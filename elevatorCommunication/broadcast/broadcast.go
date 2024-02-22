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

	udpAddr := broadcastAddr + ":" + port
	fmt.Println(udpAddr)
	
	for{
		conn, err := net.Dial("udp", udpAddr)
		if err != nil {
			fmt.Println("Failed to dial udp")
		}

		message := []byte(ip.String() + "\n")

		_, err = conn.Write(message)
		if err != nil {
			fmt.Println("Error sending UDP message: ", err)
		}

		conn.Close()
		time.Sleep(time.Millisecond * 100)
	}
}


