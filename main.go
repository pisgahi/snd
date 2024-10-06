package main

import (
	"log"
	"sync"

	"github.com/concernum/snd/client"
	"github.com/concernum/snd/server"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		server.CreateServer()
	}()

	wg.Wait()

	c := &client.Client{}

	err := c.Connect("0.0.0.0:4040")
	if err != nil {
		log.Println("Error connecting to peer:", err)
		return
	}
	defer c.Close()

	err = c.SendFile("client.go")
	if err != nil {
		log.Println("Error sending file:", err)
		return
	}

	log.Println("File sent successfully!")
}
