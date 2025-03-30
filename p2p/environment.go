package p2p

import "github.com/multiformats/go-multiaddr"

const ServerUrl = "https://pear-programming.fly.dev/api"
const ProtocolID = "/pear-programming/1.0.0"

var relayMultiAddrs = [...]string{"/dns4/pear-programming.fly.dev/tcp/3000/p2p/12D3KooWQowZAkJL3HX61Mar8hGi1eJdEoe7ybnu3mFun7aSL3NW"}

func RelayMultiAddrs() []multiaddr.Multiaddr {
	addrs := []multiaddr.Multiaddr{}
	for _, addr := range relayMultiAddrs {
		multiAddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			panic(err)
		}
		addrs = append(addrs, multiAddr)
	}
	return addrs
}
