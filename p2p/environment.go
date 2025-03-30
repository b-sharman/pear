package p2p

import "github.com/multiformats/go-multiaddr"

const ServerUrl = "http://localhost:8080/api"
const ProtocolID = "/pear-programming/1.0.0"

var relayMultiAddrs = [...]string{"/dns4/localhost/tcp/3000/p2p/12D3KooWM5NWsHa8uQ11K1FUbo1jbcrd2ceDx1Nt3H3sUcthAqMM"}

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
