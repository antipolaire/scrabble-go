package network

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"os"
	"testing"
)

func main() {

}

func TestNetwork(t *testing.T) {
	// Test the ParseFlags function
	t.Run("Smoke Test", func(t *testing.T) {
		help := flag.Bool("help", false, "Display Help")
		cfg := ParseFlags()

		if *help {
			fmt.Printf("Simple example for peer discovery using mDNS. mDNS is great when you have multiple peers in local LAN.")
			fmt.Printf("Usage: \n   Run './chat-with-mdns'\nor Run './chat-with-mdns -host [host] -port [port] -rendezvous [string] -pid [proto ID]'\n")

			os.Exit(0)
		}

		fmt.Printf("[*] Listening on: %s with port: %d\n", cfg.listenHost, cfg.listenPort)

		ctx := context.Background()
		r := rand.Reader

		// Creates a new RSA key pair for this host.
		prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			panic(err)
		}

		// 0.0.0.0 will listen on any interface device.
		sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.listenHost, cfg.listenPort))

		// libp2p.New constructs a new libp2p Host.
		// Other options can be added here.
		host, err := libp2p.New(
			libp2p.ListenAddrs(sourceMultiAddr),
			libp2p.Identity(prvKey),
		)
		if err != nil {
			panic(err)
		}

		// Set a function as stream handler.
		// This function is called when a peer initiates a connection and starts a stream with this peer.
		host.SetStreamHandler(protocol.ID(cfg.ProtocolID), handleStream)

		fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.listenHost, cfg.listenPort, host.ID())

		peerChan := initMDNS(host, cfg.RendezvousString)
		for { // allows multiple peers to join
			peer := <-peerChan // will block until we discover a peer
			fmt.Println("Found peer:", peer, ", connecting")

			if err := host.Connect(ctx, peer); err != nil {
				fmt.Println("Connection failed:", err)
				continue
			}

			// open a stream, this stream will be handled by handleStream other end
			stream, err := host.NewStream(ctx, peer.ID, protocol.ID(cfg.ProtocolID))

			if err != nil {
				fmt.Println("Stream open failed", err)
			} else {
				rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

				go writeData(rw)
				go readData(rw)
				fmt.Println("Connected to:", peer)
			}
		}
	})

}
