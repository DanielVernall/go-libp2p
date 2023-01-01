package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/event"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
)

func New() (host.Host, error) {
	//AWS VM relay @ 3.10.59.121
	knownPeerInfo, err := peer.AddrInfoFromString("/ip4/3.10.59.121/udp/4001/quic/p2p/12D3KooWGHcjc3ct1ZpKuVT4CjeuCHqFEz2tJyD9qRuCQRUev4yw")
	if err != nil {
		return nil, err
	}

	// Create a host that will attempt NAT Traversal automatically
	h, err := libp2p.New(
		libp2p.ForceReachabilityPrivate(),
		libp2p.EnableAutoRelay(
			autorelay.WithStaticRelays([]peer.AddrInfo{
				*knownPeerInfo,
			}),
			autorelay.WithNumRelays(1),
		),
	)
	if err != nil {
		return nil, err
	}

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

	log.Printf("Initial listen addresses: %v", h.Addrs())
	log.Printf("Known peer info: %v", knownPeerInfo)

	if err := h.Connect(context.Background(), *knownPeerInfo); err != nil {
		log.Printf("Could not connect to known peer: %v", err)
	}

	log.Printf("Connected to known peer")

	_, err = h.NewStream(context.Background(), knownPeerInfo.ID, "/ipfs/ping/1.0.0")
	if err != nil {
		log.Printf("Ping failed: %v", err)
	}

	// Manual relay - shouldn't be needed....
	// r, err := client.Reserve(context.Background(), h, *knownPeerInfo)
	// if err != nil {
	// 	log.Printf("Failed to make relay reservation: %v", err)
	// }

	// log.Printf("Made relay reservation: %v", r)
	log.Printf("Updated listen addresses: %v", h.Addrs())

	return h, nil
}

func main() {
	_, err := New()
	if err != nil {
		log.Printf("Error creating host %v", err)
	}

	select {} // Wait forever
}
