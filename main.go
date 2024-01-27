package main

// TODO(alx): Give a try to Cobra.

import (
	"flag"
	"strings"
)

func main() {
	options := Options{}

	flag.Var(&options.Ports, "port", "Ports to connect to")
	flag.StringVar(&options.Network, "network", "tcp", "Network protocol [tcp|udp]")
	flag.StringVar(&options.Address, "address", "localhost", "Address of the server to connect to.")

	help := flag.Bool("help", false, "Display all commands.")
	instance := flag.String("instance", "server", "Boots up either a client or a server.")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	var linstance = strings.ToLower(*instance)
	if linstance == "server" {
		server = NewServer()
		server.Run(&options)
	} else if linstance == "client" {
		Run(&options)
	}
}
