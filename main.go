package main

import (
	"github.com/raoulh/acr122u"
)

func main() {

	// Discover all services on the network
	discoverBinkyServer()

	/*ctx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}

	h := &handler{log.New(os.Stdout, "", 0)}

	ctx.Serve(h)
	*/
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
