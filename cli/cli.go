package sndCli

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	StartServer bool
	FileToSend  string
	ServerAddr  string
	Terminate   bool
}

func ParseFlags() *Flags {
	var flags Flags

	flag.BoolVar(&flags.StartServer, "s", false, "Start server")
	flag.StringVar(&flags.ServerAddr, "to", "", "Recipient")
	flag.StringVar(&flags.FileToSend, "file", "", "File to send")
	flag.BoolVar(&flags.Terminate, "t", false, "Terminate server")

	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Snd is a file transfer program utilizing TCP.")
		os.Exit(0)
	}

	return &flags
}
