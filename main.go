package main

import (
	"log"
	"sync"

	certgen "github.com/concernum/snd/cert"
	sndCli "github.com/concernum/snd/cli"
	"github.com/concernum/snd/client"
	"github.com/concernum/snd/server"
	"github.com/concernum/snd/sndcfg"
)

func main() {
	flags := sndCli.ParseFlags()

	config := sndcfg.LoadOrCreateConfig(flags)

	certgen.SetupCertificates()

	var wg sync.WaitGroup

	if flags.StartServer {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Starting server...")
			server.CreateServer(config.ServerAddr, config.ReceivedDir)
		}()
	}

	if flags.FileToSend != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := &client.Client{}
			client.HandleFileSending(c, flags, config)
		}()
	}

	wg.Wait()
}
