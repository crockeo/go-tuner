package server

import (
	"fmt"
	"github.com/crockeo/go-tuner/synth"
	"net"
	"strings"
	"time"
)

// Handling a particular TCP connection.
func handleTCPConnection(conn *net.TCPConn, noteChannel chan synth.DelayedNoteData) {
	defer conn.Close()
	buffer := make([]byte, 512)
	for {
		rlen, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read from TCP socket: " + err.Error())
			return
		}

		if rlen > 0 {
			strs := strings.Split(strings.TrimSpace(string(buffer[:2])), "\n")
			for _, v := range strs {
				err = HandleMessage(v, noteChannel)
				if err != nil {
					fmt.Println("Failed to handle message \"" + v + "\": " + err.Error())
				}
			}
		}

		buffer = make([]byte, 64)
		time.Sleep(10 * time.Millisecond)
	}
}

// Starting the go-tuner server either on the main thread or another thread.
func Start(errChannel chan error, noteChannel chan synth.DelayedNoteData) {
	addr := net.TCPAddr{
		Port: 3000,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// Creating a connection.
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		errChannel <- err
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			errChannel <- err
			continue
		}

		go handleTCPConnection(conn, noteChannel)
	}
}
