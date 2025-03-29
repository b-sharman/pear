package p2p

import "github.com/multiformats/go-multiaddr"

const ServerUrl = "https://pear-programming.fly.dev/api"
const ProtocolID = "/pear-programming/1.0.0"
const RelayPeerID = "12D3KooWDpUzpsoB1nGF8T31iouGcEiCFNKn2Zezto3kEfKECxEw"

var relayMultiAddrs = [...]string{"/dns4/pear-programming.fly.dev/tcp/1337", "/dns6/pear-programming.fly.dev/tcp/1337"}

func RelayMultiAddrs() []multiaddr.Multiaddr {
	addrs := []multiaddr.Multiaddr{}
	for _, addr := range relayMultiAddrs {
		multiAddr, _ := multiaddr.NewMultiaddr(addr)
		addrs = append(addrs, multiAddr)
	}
	return addrs
}
