package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	conn net.Conn
}

const ChunkSize = 64 * 1024

func (c *Client) Connect(address string) error {
	var err error
	for {
		c.conn, err = net.Dial("tcp", address)
		if err == nil {
			log.Println("Connected to peer:", address)
			return nil
		}
		log.Println("Error connecting to peer:", err)
		log.Println("Retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}
}

func (c *Client) SendFile(fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", fileName)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	_, err = c.conn.Write([]byte(fmt.Sprintf("%s:%d\n", fileName, fileSize)))
	if err != nil {
		return err
	}

	buffer := make([]byte, ChunkSize)
	chunkIndex := 0
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if bytesRead == 0 {
			break
		}

		chunkHeader := fmt.Sprintf("Chunk %d of %d\n", chunkIndex, (fileSize+ChunkSize-1)/ChunkSize)
		_, err = c.conn.Write([]byte(chunkHeader))
		if err != nil {
			return err
		}

		_, err = c.conn.Write(buffer[:bytesRead])
		if err != nil {
			return err
		}

		log.Printf("Sent chunk %d (%d bytes)\n", chunkIndex, bytesRead)
		chunkIndex++
	}

	log.Println("File sent successfully in chunks")
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
		log.Println("Connection closed")
	}
	return nil
}
