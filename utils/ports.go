package utils

import (
	"fmt"
	"math/rand"
	"net"
)

// Attempts to connect to the port via TCP. If something accepts the connection
// then IsPortFree returns false. If nothing accepts the connection it will
// return true.
func IsPortFree(port int) bool {
	local, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if conn, err := net.DialTCP("tcp", nil, local); err == nil {
		defer func() { conn.Close() }()
		return false
	}
	return true
}

// Will asynchronously discover an available port. To retrieve the port number
// you need to consume from the returned channel.
func FindAvailablePort() <-chan int {
	portChannel := make(chan int, 1)
	go func() {
		port, free := randomPortInfo()
		for !free {
			port, free = randomPortInfo()
		}
		portChannel <- port
	}()
	return portChannel
}

func randomPortInfo() (port int, free bool) {
	port = 3000 + rand.Intn(2000)
	free = IsPortFree(port)
	return
}
