package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func CreateServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error starting TCP Server:", err)
		return
	}
	defer listener.Close()

	log.Println("TCP Server started")

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

	reader := bufio.NewReader(conn)

	fileName, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading filename:", err)
		return
	}
	fileName = strings.TrimSpace(fileName)

	receivedFile, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer receivedFile.Close()

	buffer := make([]byte, 1024) // 1 KB buffer size
	for {
		bytesRead, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("File reception completed from:", conn.RemoteAddr())
				break
			}
			log.Println("Error reading from connection:", err)
			return
		}

		_, err = receivedFile.Write(buffer[:bytesRead])
		if err != nil {
			log.Println("Error writing to file:", err)
			return
		}
	}

	log.Println("File received successfully from:", conn.RemoteAddr())
}
