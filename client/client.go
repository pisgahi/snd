package client

import (
	"log"
	"net"
)

func ConnectToPeer() {
	conn, err := net.Dial("tcp", "localhost:4040")
	if err != nil {
		log.Println("Error connecting to peer:", err)
		return
	}

	defer conn.Close()

	log.Print("\nConnected to peer")
}
