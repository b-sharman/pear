package client

import (
	"context"
	"fmt"
	"github.com/b-sharman/pear/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
	"golang.org/x/term"
	"io"
	"os"
)

func Start(ctx context.Context, roomid string) {
	relay, err := peer.AddrInfoFromP2pAddr(p2p.RelayMultiAddrs()[0])
	if err != nil {
		panic(err)
	}

	fmt.Println("relay", relay)

	// start a libp2p node that listens on a random local TCP port
	node, err := libp2p.New(
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
		libp2p.EnableAutoRelayWithStaticRelays([]peer.AddrInfo{*relay}),
	)
	if err != nil {
		panic(err)
	}

	// connect to the relay
	if err := node.Connect(context.Background(), *relay); err != nil {
		panic(err)
	}
	fmt.Println("Connected to relay!")

	// Print this node's `PeerInfo` in multiaddr format
	peerInfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address:", addrs[0])

	addrStr, err := getAddr(roomid)
	if err != nil {
		fmt.Printf("encountered err on getting addr: %s\n", err.Error())
		return
	}
	fmt.Printf("connecting to %s\n", addrStr)

	// Get a peer.AddrInfo from addrStr
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}

	relayaddr, err := multiaddr.NewMultiaddr("/p2p/" + relay.ID.String() + "/p2p-circuit/p2p/" + peerAddrInfo.ID.String())
	peerrelayinfo := peer.AddrInfo{
		ID:    peerAddrInfo.ID,
		Addrs: []multiaddr.Multiaddr{relayaddr},
	}

	// Connect to the peer node
	if err := node.Connect(ctx, peerrelayinfo); err != nil {
		panic(err)
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	stream, err := node.NewStream(network.WithAllowLimitedConn(ctx, "connect"), peerrelayinfo.ID, "/connect/0.0.0")
	if err != nil {
		panic(err)
	}

	go func() { _, _ = io.Copy(os.Stdout, stream) }()
	go func() { _, _ = io.Copy(stream, os.Stdin) }()
	<-ctx.Done()
	term.Restore(int(os.Stdin.Fd()), oldState)

	if err := stream.Close(); err != nil {
		panic(err)
	}

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}
}
