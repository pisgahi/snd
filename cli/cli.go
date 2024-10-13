package sndCli

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	FileToSend  string
	ReceivedDir string
	CertDir     string
	StartServer bool
	Terminate   bool
	To          string
	ServerAddr  string
	Help        bool
}

func ParseFlags() *Flags {
	flags := &Flags{}

	flag.StringVar(&flags.FileToSend, "f", "", "File to send")
	flag.StringVar(&flags.ReceivedDir, "dir", "", "Directory for received files")
	flag.StringVar(&flags.CertDir, "cert", "", "Directory for certificates")
	flag.BoolVar(&flags.StartServer, "s", false, "Start server")
	flag.BoolVar(&flags.Terminate, "t", false, "Terminate server")
	flag.StringVar(&flags.To, "to", "", "Recipient address")
	flag.StringVar(&flags.ServerAddr, "addr", "", "Server address")
	flag.BoolVar(&flags.Help, "h", false, "Display help information")

	flag.Parse()

	if flags.Help {
		printHelp()
		os.Exit(0)
	}

	if flag.NFlag() == 0 {
		fmt.Println("Snd is a file transfer program utilizing TCP.")
		os.Exit(0)
	}

	return flags
}

func printHelp() {
	fmt.Println("Snd is a file transfer program utilizing TCP.")
	flag.PrintDefaults()
}
