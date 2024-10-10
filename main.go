package main

import (
	"log"
	"sync"

	sndCli "github.com/concernum/snd/cli"
	"github.com/concernum/snd/client"
	"github.com/concernum/snd/server"
)

func main() {
	flags := sndCli.ParseFlags()

	var wg sync.WaitGroup

	if flags.StartServer {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Starting server...")
			server.CreateServer("0.0.0.0:4040")
		}()
	}

	if flags.FileToSend != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			c := &client.Client{}
			err := c.Connect(flags.ServerAddr)
			if err != nil {
				log.Println("Error connecting to peer:", err)
				return
			}

			err = c.SendFile(flags.FileToSend)
			if err != nil {
				log.Println("Error sending file:", err)
				return
			}

			log.Println("File sent successfully.")

			if flags.Terminate {
				log.Println("Terminating connection per -t flag.")
			} else {
				log.Println("Still connected...")
				select {}
			}

			err = c.Close()
			if err != nil {
				log.Println("Error closing connection:", err)
			} else {
				log.Println("Connection closed.")
			}
		}()
	}

	wg.Wait()
}
