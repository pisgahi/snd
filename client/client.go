package client

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	conn net.Conn
}

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

	_, err = c.conn.Write([]byte(fmt.Sprintf("%s\n", fileName))) // Sending filename followed by newline
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024) // 1 KB buffer size
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("File transmission completed")
				break
			}
			return err
		}

		_, err = c.conn.Write(buffer[:bytesRead])
		if err != nil {
			return err
		}
	}

	log.Println("File sent successfully")
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
