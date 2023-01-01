package main

import (
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/event"
)

func main() {
	//run()
	test()
}

func run() {
	addrs := []string{
		"/ip4/0.0.0.0/udp/4001/quic",
		"/ip6/::/udp/4001/quic",
		"/ip4/0.0.0.0/tcp/4001",
		"/ip6/::/tcp/4001",
	}

	h, err := libp2p.New(
		// libp2p.EnableNATService(),
		libp2p.DisableRelay(),
		libp2p.EnableRelayService(),
		libp2p.ForceReachabilityPublic(),
		libp2p.ListenAddrStrings(addrs...),
	)
	if err != nil {
		log.Panicf("Error creating host: %v", err)
	}

	subidentified, err := h.EventBus().Subscribe(&event.EvtPeerIdentificationCompleted{})
	if err != nil {
		log.Panicf("Error subscribing to PeerIdentificationCompleted event: %v", err)
	}
	defer subidentified.Close()

	log.Printf("Subscribed to PeerIdentificationCompleted event")
	go func() {
		for change := range subidentified.Out() {
			log.Printf("Identified: %v", change)
		}
	}()

	subconnectedness, err := h.EventBus().Subscribe(&event.EvtPeerConnectednessChanged{})
	if err != nil {
		log.Panicf("Error subscribing to PeerConnectedness event: %v", err)
	}
	defer subconnectedness.Close()

	log.Printf("Subscribed to PeerConnectedness event")
	go func() {
		for change := range subconnectedness.Out() {
			log.Printf("Connectedness change: %v", change)
		}
	}()

	sublocalprotocol, err := h.EventBus().Subscribe(&event.EvtLocalProtocolsUpdated{})
	if err != nil {
		log.Panicf("Error subscribing to LocalProtocolsUpdated event: %v", err)
	}
	defer sublocalprotocol.Close()

	log.Printf("Subscribed to LocalProtocolsUpdated event")
	go func() {
		for change := range sublocalprotocol.Out() {
			log.Printf("LocalProtocols updated: %v", change)
		}
	}()

	sublocalreach, err := h.EventBus().Subscribe(&event.EvtLocalReachabilityChanged{})
	if err != nil {
		log.Panicf("Error subscribing to LocalReachabilityChanged event: %v", err)
	}
	defer sublocalreach.Close()

	log.Printf("Subscribed to LocalReachabilityChanged event")
	go func() {
		for change := range sublocalreach.Out() {
			log.Printf("LocalReachability changed: %v", change)
		}
	}()

	log.Printf("Started listening on: %v", h.Addrs())
	log.Printf("Peer ID: %v", h.ID().String())
	log.Printf("Relay service running...")
	log.Printf("NAT service running...")

	log.Printf("Registered protocols: %v", h.Mux().Protocols())

	select {} // wait forever
}
