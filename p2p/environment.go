package p2p

import "github.com/multiformats/go-multiaddr"

const ServerUrl = "https://pear-programming.fly.dev/api"

var relayMultiAddrs = [...]string{"/dns4/pear-programming.fly.dev/tcp/3000/p2p/12D3KooWM5NWsHa8uQ11K1FUbo1jbcrd2ceDx1Nt3H3sUcthAqMM"}

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
