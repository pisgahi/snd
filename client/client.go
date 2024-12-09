package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	sndCli "github.com/pisgahi/snd/cli"
	"github.com/pisgahi/snd/sndcfg"
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
			fmt.Println("Connected to peer:", address)
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

		fmt.Printf("Sent chunk %d (%d bytes)\n", chunkIndex, bytesRead)
		chunkIndex++
	}

	fmt.Println("File sent successfully in chunks")
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
		fmt.Println("Connection closed")
	}
	return nil
}

func HandleFileSending(c *Client, flags *sndCli.Flags, config *sndcfg.Config) {
	addr := flags.ServerAddr
	if addr == "" {
		addr = config.ServerAddr
	}

	if err := c.Connect(addr); err != nil {
		log.Println("Error connecting to peer:", err)
		return
	}

	if err := c.SendFile(flags.FileToSend); err != nil {
		log.Println("Error sending file:", err)
		return
	}

	fmt.Println("File sent successfully.")

	if flags.Terminate {
		fmt.Println("Terminating connection per -t flag.")
	} else {
		fmt.Println("Still connected...")
		select {}
	}

	if err := c.Close(); err != nil {
		log.Println("Error closing connection:", err)
	} else {
		fmt.Println("Connection closed.")
	}
}
