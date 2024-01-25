package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen to incoming udp packets
	addr := getIP()

	// Send and receive UDP on different port
	// sendUDP(addr)

	// receiveUDP()

	// Send and receive TCP

	connection, err := net.Dial("tcp", addr+":34933")
	if err != nil {
		panic(err)
	}

	defer connection.Close()

	buffer := make([]byte, 1024)
	_, err = connection.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Received: ", string(buffer))

	// Accept incoming connections
	serverAddr, err := net.ResolveTCPAddr("tcp", "10.100.23.24:12345")
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	_, err = connection.Write([]byte("Connect to: 10.100.23.24:12345\000"))
	if err != nil {
		panic(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}

	fmt.Print("Accepted connection from: ", conn.RemoteAddr())

	buffer = make([]byte, 1024)
	_, err = connection.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Received: ", string(buffer))

}

func getIP() string {
	serverAddr, err := net.ResolveUDPAddr("udp", ":30000")
	if err != nil {
		panic(err)
	}

	connection, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	buffer := make([]byte, 1024)
	n, _, err := connection.ReadFromUDP(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Recieved: ", string(buffer))

	return string(buffer[25 : n-1])
}

func sendUDP(addr string) {
	connection1, err := net.Dial("udp", addr+":20014")
	if err != nil {
		panic(err)
	}
	defer connection1.Close()

	_, err = connection1.Write([]byte("Don't @ me bro!\000"))
	if err != nil {
		panic(err)
	}
}

func receiveUDP() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":20014")
	if err != nil {
		panic(err)
	}

	connection, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
}
