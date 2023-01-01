package main

import (
	"log"

	"github.com/libp2p/go-libp2p"
)

func main() {
	h, err := libp2p.New(
		libp2p.EnableNATService(),
		libp2p.EnableRelayService(),
	)

	if err != nil {
		log.Printf("Error creating host: %v", err)
	}

	log.Printf("ID: %v", h.ID())
	log.Printf("Listening on: %v", h.Network().ListenAddresses())
	log.Printf("Handling protocols: %v", h.Mux().Protocols())
}
