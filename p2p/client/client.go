package client

import (
	// "bufio"
	"context"
	"fmt"
	"io"
	// "os"

	"github.com/b-sharman/pear/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	multiaddr "github.com/multiformats/go-multiaddr"
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
		libp2p.EnableAutoRelayWithStaticRelays([]peerstore.AddrInfo{*relay}),
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
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
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

	// Get a peerstore.AddrInfo from addrStr
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}

	relayaddr, err := multiaddr.NewMultiaddr("/p2p/" + relay.ID.String() + "/p2p-circuit/p2p/" + peer.ID.String())
	peerrelayinfo := peerstore.AddrInfo{
		ID:    peer.ID,
		Addrs: []multiaddr.Multiaddr{relayaddr},
	}

	// Connect to the peer node
	if err := node.Connect(ctx, peerrelayinfo); err != nil {
		panic(err)
	}

	stream, err := node.NewStream(ctx, peerrelayinfo.ID, protocol.ID(p2p.ProtocolID))

	// reader := bufio.NewReader(stream)
	// writer := bufio.NewWriter(stream)

	// go func() {
	buf := make([]byte, 1024)

	for {
		n, err := stream.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			fmt.Print(buf[:n])
		}
	}
	// _, err = io.Copy(os.Stdout, reader)
	// }()

	// select {
	// case <-ctx.Done():
	// 	return
	// default:
	// 	// actual long work
	// 	io.Copy(rw.Writer, os.Stdin)
	// }

	// stream.Close()

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}
}
