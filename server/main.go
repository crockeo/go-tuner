package server

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Handling a particular TCP connection.
func handleTCPConnection(conn *net.TCPConn) {
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
				err = HandleMessage(v)
				if err != nil {
					fmt.Println("Failed to handle message \"" + v + "\": " + err.Error())
				}
			}
		}

		buffer = make([]byte, 64)
		time.Sleep(10 * time.Millisecond)
	}
}

// The backend function to start the server.
func realStart(handleErrors bool) error {
	addr := net.TCPAddr{
		Port: 3000,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// Creating a connection.
	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		if !handleErrors {
			return err
		}

		// Handling the error itself.
		fmt.Println("Failed to open TCP listener: " + err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			if !handleErrors {
				return err
			}

			fmt.Println("Failed to accept TCP connection: " + err.Error())
			os.Exit(1)
		}

		go handleTCPConnection(conn)
	}

	return nil
}

// Starting the go-tuner server either on the main thread or another thread.
func Start(threaded bool) error {
	if threaded {
		go realStart(true)
	} else {
		err := realStart(false)
		if err != nil {
			return err
		}
	}

	return nil
}
