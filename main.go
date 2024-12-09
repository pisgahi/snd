package main

import (
	"fmt"
	"sync"

	certgen "github.com/pisgahi/snd/cert"
	sndCli "github.com/pisgahi/snd/cli"
	"github.com/pisgahi/snd/client"
	"github.com/pisgahi/snd/server"
	"github.com/pisgahi/snd/sndcfg"
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
			fmt.Println("Starting server...")
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
