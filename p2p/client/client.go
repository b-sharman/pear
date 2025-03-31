package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/b-sharman/pear/p2p"
	"github.com/creack/pty"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	multiaddr "github.com/multiformats/go-multiaddr"
	"golang.org/x/term"
)

func Start(ctx context.Context, roomid string, username string) {
	relay, err := peer.AddrInfoFromP2pAddr(p2p.RelayMultiAddrs()[0])
	if err != nil {
		panic(err)
	}

	fmt.Println("Relay", relay.ID)

	node, err := libp2p.New(
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
		libp2p.EnableAutoRelayWithStaticRelays([]peer.AddrInfo{*relay}),
	)
	if err != nil {
		panic(err)
	}

	// connect to the relay
	if err := node.Connect(ctx, *relay); err != nil {
		panic(err)
	}
	fmt.Println("Connected to relay!")

	addrStr, err := getAddr(roomid)
	if err != nil {
		fmt.Printf("encountered err on getting addr: %s\n", err.Error())
		return
	}
	fmt.Printf("Connecting to %s\n", addrStr)

	// Get a peer.AddrInfo from addrStr
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}
	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}

	peerAddr, err := multiaddr.NewMultiaddr("/p2p/" + relay.ID.String() + "/p2p-circuit/p2p/" + peerAddrInfo.ID.String())
	peerInfo := peer.AddrInfo{
		ID:    peerAddrInfo.ID,
		Addrs: []multiaddr.Multiaddr{peerAddr},
	}

	ctx = network.WithAllowLimitedConn(ctx, "REASON")
	bundle := network.NotifyBundle{
		ConnectedF: func(_ network.Network, conn network.Conn) {
			if conn.RemotePeer() == peerInfo.ID {
				fmt.Println("Connected to " + conn.RemoteMultiaddr().String() + " ID: " + conn.ID())

				sendUsername(username, ctx, node, peerInfo.ID)
				streamStdIO(ctx, node, peerInfo.ID)
				resizeHandler(ctx, node, peerInfo.ID)
			}
		},
		DisconnectedF: func(_ network.Network, conn network.Conn) {
			fmt.Println("Disconnected from " + conn.RemoteMultiaddr().String() + " ID: " + conn.ID())
		},
	}
	node.Network().Notify(&bundle)

	// Connect to the peer node
	if err := node.Connect(ctx, peerInfo); err != nil {
		panic(err)
	}
	<-ctx.Done()

	// shut the node down
	defer func() {
		if err := node.Close(); err != nil {
			panic(err)
		}
	}()
}

func sendUsername(username string, ctx context.Context, node host.Host, peerid peer.ID) {
	stream, err := node.NewStream(ctx, peerid, "/username/0.0.0")
	if err != nil {
		panic(err)
	}

	stream.Write([]byte(username))
	stream.Write([]byte{'\n'})
}

func resizeHandler(ctx context.Context, node host.Host, peerid peer.ID) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	stream, err := node.NewStream(ctx, peerid, "/resize/0.0.0")
	if err != nil {
		panic(err)
	}

	go func() {
		for range ch {
			size, err := pty.GetsizeFull(os.Stdout)
			if err != nil {
				panic("Unable to read terminal size")
			}

			marshaledSize, err := json.Marshal(size)
			if err != nil {
				panic("Unable to marshal size")
			}

			stream.Write(marshaledSize)
			stream.Write([]byte{'\n'})
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.
}

func streamStdIO(ctx context.Context, node host.Host, peerid peer.ID) {
	stream, err := node.NewStream(ctx, peerid, "/connect/0.0.0")
	if err != nil {
		panic(err)
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	go func() { _, _ = io.Copy(os.Stdout, stream) }()
	_, _ = io.Copy(stream, os.Stdin)

	stream.Reset()
}
