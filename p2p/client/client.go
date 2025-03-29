package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/b-sharman/pear/p2p"
	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func Start(ctx context.Context, roomid string) {
	// start a libp2p node that listens on a random local TCP port
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	// print the node's PeerInfo in multiaddr format
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

	// get a peerstore.AddrInfo from addrStr
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}

	// connect to the peer node
	if err := node.Connect(ctx, *peer); err != nil {
		panic(err)
	}

	stream, err := node.NewStream(ctx, peer.ID, protocol.ID(p2p.ProtocolID))

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go func() {
		io.Copy(os.Stdout, rw.Reader)
	}()

	select {
	case <-ctx.Done():
		return
	default:
		// actual long work
		io.Copy(rw.Writer, os.Stdin)
	}

	stream.Close()

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}
}
