package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const ChunkSize = 64 * 1024 // 64 KB

func CreateServer(address, baseDir string) {
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

		go handleConnection(conn, baseDir)
	}
}

func handleConnection(conn net.Conn, baseDir string) {
	defer conn.Close()
	log.Println("Handling connection from:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	fileMeta, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading file metadata:", err)
		return
	}
	fileMeta = strings.TrimSpace(fileMeta)
	parts := strings.Split(fileMeta, ":")
	if len(parts) != 2 {
		log.Println("Invalid file metadata format:", fileMeta)
		return
	}
	fileName := parts[0]
	fileSize, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		log.Println("Invalid file size:", parts[1])
		return
	}
	log.Printf("Receiving file: %s (%d bytes)\n", fileName, fileSize)

	err = os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		log.Println("Error creating base directory:", err)
		return
	}

	escapedFilePath := filepath.Join(baseDir, filepath.Base(fileName))

	receivedFile, err := os.Create(escapedFilePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer receivedFile.Close()

	var totalBytesReceived int64
	for {
		chunkHeader, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("File reception completed from:", conn.RemoteAddr())
				break
			}
			log.Println("Error reading chunk header:", err)
			return
		}
		log.Printf("Chunk header received: %s\n", strings.TrimSpace(chunkHeader))

		bytesToRead := int64(ChunkSize)
		if totalBytesReceived+ChunkSize > fileSize {
			bytesToRead = fileSize - totalBytesReceived
		}

		buffer := make([]byte, bytesToRead)
		bytesRead, err := io.ReadFull(reader, buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("End of file detected")
				break
			}
			log.Println("Error reading chunk data:", err)
			return
		}

		_, err = receivedFile.Write(buffer[:bytesRead])
		if err != nil {
			log.Println("Error writing to file:", err)
			return
		}

		totalBytesReceived += int64(bytesRead)
		log.Printf("Received chunk (%d bytes), total received: %d/%d bytes\n", bytesRead, totalBytesReceived, fileSize)

		if totalBytesReceived >= fileSize {
			log.Println("File reception completed")
			break
		}
	}

	log.Println("File received successfully and saved to:", escapedFilePath)
}
