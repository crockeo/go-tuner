package server

import (
	"fmt"
	"net"
	"os"
	"time"
)

func realStart(handleErrors bool) error {
	addr := net.UDPAddr{
		Port: 3000,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// Creating a connection.
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		if !handleErrors {
			return err
		}

		// Handling the error itself.
		fmt.Println("Failed to open UDP listener: " + err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	var buffer []byte
	for {
		rlen, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read buffer information: " + err.Error())
			continue
		}

		if rlen > 0 {
			fmt.Println(string(buffer))
		}

		time.Sleep(32 * time.Millisecond)
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
