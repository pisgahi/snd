package main

import (
	"flag"
	"log"
	"sync"

	"github.com/concernum/snd/client"
	"github.com/concernum/snd/server"
)

func main() {
	startServer := flag.Bool("start", false, "Start Server")
	serverAddr := flag.String("usr", "", "Recipient")
	fileToSend := flag.String("file", "", "File to send to usr")
	terminate := flag.Bool("t", false, "Terminate server")

	flag.Parse()

	var wg sync.WaitGroup

	if *startServer {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Starting server...")
			server.CreateServer("0.0.0.0:4040")
		}()
	}

	if *fileToSend != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()

			c := &client.Client{}
			err := c.Connect(*serverAddr)
			if err != nil {
				log.Println("Error connecting to peer:", err)
				return
			}

			err = c.SendFile(*fileToSend)
			if err != nil {
				log.Println("Error sending file:", err)
				return
			}

			log.Println("File sent successfully.")

			if *terminate {
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
