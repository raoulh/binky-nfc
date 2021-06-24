package main

import (
	"fmt"
	"log"
	"os"

	"github.com/raoulh/acr122u"
	"github.com/raoulh/binky-nfc/mdns"
)

func main() {

	// Discover all services on the network
	for {
		host, err := mdns.DiscoverBinkyServer()
		if err == nil {
			fmt.Println("Found host:", host)
			break
		}
	}

	ctx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}

	h := &handler{log.New(os.Stdout, "", 0)}

	ctx.Serve(h)
}

type handler struct {
	acr122u.Logger
}

func (h *handler) ServeCard(c acr122u.Card) {
	h.Printf("%x\n", c.UID())
}

func (h *handler) CardRemoved() {
	h.Printf("card removed\n")
}
