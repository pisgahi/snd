package sndCli

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	FileToSend  string
	ReceivedDir string
	StartServer bool
	Terminate   bool
	To          string
	ServerAddr  string
}

func ParseFlags() *Flags {
	flags := &Flags{}

	flag.StringVar(&flags.FileToSend, "f", "", "File to send")
	flag.StringVar(&flags.ReceivedDir, "dir", "", "Directory for received files")
	flag.BoolVar(&flags.StartServer, "s", false, "Start server")
	flag.BoolVar(&flags.Terminate, "t", false, "Terminate server")
	flag.StringVar(&flags.To, "to", "", "Recipient")
	flag.StringVar(&flags.ServerAddr, "addr", "", "Server address")

	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Snd is a file transfer program utilizing TCP.")
		os.Exit(0)
	}

	return flags
}
