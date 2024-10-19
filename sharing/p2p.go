package main

import (
	"context"
	"fmt"
	"log"
	"os"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	network "github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// Create a new libp2p host
func createHost(port string) (host.Host, error) {
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", port),
		),
	)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Host created. We are: %s\n", h.ID())
	fmt.Println(h.Addrs())
	return h, nil
}

// Stream handler for incoming connections
func handleStream(s network.Stream) {
	fmt.Println("Got a new stream!")
	// Placeholder for stream communication logic (read/write data)
}

// Function to connect to a peer using their multiaddress
func connectToPeer(h host.Host, peerAddr string) error {
	addr, err := multiaddr.NewMultiaddr(peerAddr)
	if err != nil {
		return err
	}

	info, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return err
	}

	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	s, err := h.NewStream(context.Background(), info.ID, "/myapp/1.0.0")
	if err != nil {
		return err
	}

	fmt.Println("Connected to", info.ID)
	handleStream(s)

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a port number")
	}

	port := os.Args[1]
	h, err := createHost(port)
	if err != nil {
		log.Fatal(err)
	}

	// Set a handler for incoming streams
	h.SetStreamHandler("/myapp/1.0.0", handleStream)

	if len(os.Args) > 2 {
		peerAddr := os.Args[2]
		fmt.Printf("Connecting to %s...\n", peerAddr)
		if err := connectToPeer(h, peerAddr); err != nil {
			log.Fatal(err)
		}
	}

	// Keep the host alive
	select {}
}
