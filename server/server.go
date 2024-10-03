package server

import (
	"log"
	"net"

	"github.com/concernum/snd/client"
)

func CreateServer() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Println("Error starting TCP Server:", err)
		return
	}

	defer listener.Close()

	log.Println("TCP Server started")

	go client.ConnectToPeer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Handling connection from:", conn.RemoteAddr())
}
