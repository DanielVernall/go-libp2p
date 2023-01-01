package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/autonat"
)

func test() {
	// Nodes that provide the AutoNAT service are preconfigured to help other peers
	// to determine their NAT status, and whether they are publically reachable.
	//createAutoNATServiceHost()

	h, err := createHostUnderTest()
	if err != nil {
		log.Printf("Failed to create host: %v", err)
	}

	// Nodes connect to the AutoNAT service provider on the AutoNAT protocol
	// to request a dialback. The AutoNAT service provider will respond with
	// whether the node is publically accessible on their observed dial address.
	err = determineNATStatus(h)
	if err != nil {
		log.Printf("Failed to determine NAT status: %v", err)
	}
}

// func createAutoNATServiceHost() {

// }

func createHostUnderTest() (host.Host, error) {
	h, err := libp2p.New()
	if err != nil {
		return nil, err
	}

	return h, nil
}

func determineNATStatus(h host.Host) error {
	// Connect to the AutoNat service
	peerinfo, err := peer.AddrInfoFromString("/ip4/35.178.193.98/tcp/4001/p2p/12D3KooWGHcjc3ct1ZpKuVT4CjeuCHqFEz2tJyD9qRuCQRUev4yw")
	if err != nil {
		return err
	}

	err = h.Connect(context.Background(), *peerinfo)
	if err != nil {
		return err
	}

	log.Printf("Connected to peer!")

	// Create autonat client and request dialback
	ac := autonat.NewAutoNATClient(h, nil)
	ma, err := ac.DialBack(context.Background(), peerinfo.ID)
	if err != nil {
		return err
	}

	log.Printf("Dialback response: %v", ma)

	return nil
}
