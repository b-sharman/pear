package host

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"

	"github.com/b-sharman/pear/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func Start(roomid string) error {
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
	err = registerRoom(addrs, roomid)
	if err != nil {
		return err
	}

	node.SetStreamHandler(p2p.ProtocolID, func(s network.Stream) {
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		cmd := exec.Command("man", "cat")
		cmd.Stdin = rw.Reader
		cmd.Stdout = rw.Writer

		cmd.Run()
	})

	return nil
}

func registerRoom(addrs []multiaddr.Multiaddr, roomid string) error {
	b := bytes.NewBufferString(addrs[0].String())

	u, _ := url.Parse(p2p.ServerUrl)
	u = u.JoinPath("register", roomid)
	resp, err := http.Post(u.String(), "", b)

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("name server returned %s\n", resp.Status)
	}
	return nil
}
