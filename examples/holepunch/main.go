package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	pb "github.com/libp2p/go-libp2p/p2p/protocol/identify/pb"
	"github.com/libp2p/go-msgio/protoio"
	ma "github.com/multiformats/go-multiaddr"
)

// Implement the hole punching procedure. https://docs.libp2p.io/concepts/hole-punching/
// Assume private instead of using AutoNAT

func main() {
	run()
}

func run() {
	h1, err := libp2p.New()
	if err != nil {
		log.Printf("Failed to create h1: %v", err)
		return
	}

	log.Printf("h1 addresses: %v", h1.Addrs())

	h2info, err := peer.AddrInfoFromString("/ip4/3.10.59.121/udp/4001/quic/p2p/12D3KooWGHcjc3ct1ZpKuVT4CjeuCHqFEz2tJyD9qRuCQRUev4yw")
	if err != nil {
		log.Printf("Failed to create h2info from string: %v", err)
		return
	}

	// h2, err := libp2p.New()
	// if err != nil {
	// 	log.Printf("Failed to create h2: %v", err)
	// 	return
	// }

	// log.Printf("h2 addresses: %v", h2.Addrs())

	// h2info := peer.AddrInfo{
	// 	ID:    h2.ID(),
	// 	Addrs: h2.Addrs(),
	// }

	if err := h1.Connect(context.Background(), *h2info); err != nil {
		log.Printf("Failed to connect h1 and h2: %v", err)
		return
	}

	s, err := h1.NewStream(context.Background(), h2info.ID, "/ipfs/id/1.0.0")
	if err != nil {
		log.Printf("h1 couldn't open identify stream with h2. %v", err)
	}

	r := protoio.NewDelimitedReader(s, 8*1024)
	mes := &pb.Identify{}

	if err := readAllIDMessages(r, mes); err != nil {
		log.Printf("error reading identify message: %v", err)
		s.Reset()
		return
	}

	log.Printf("%s received message from %s %s", s.Protocol(), s.Conn().RemotePeer(), s.Conn().RemoteMultiaddr())
	maddr, err := ma.NewMultiaddrBytes(mes.GetObservedAddr())
	if err != nil {
		log.Printf("error parsing received observed addr for %s: %s", s.Conn(), err)
		return
	}

	log.Printf("Observed Address: %v", maddr)

	defer s.Close()
}

func readAllIDMessages(r protoio.Reader, finalMsg proto.Message) error {
	mes := &pb.Identify{}
	for i := 0; i < 250; i++ {
		switch err := r.ReadMsg(mes); err {
		case io.EOF:
			return nil
		case nil:
			proto.Merge(finalMsg, mes)
		default:
			return err
		}
	}

	return fmt.Errorf("too many parts")
}
