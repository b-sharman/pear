package host

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/b-sharman/pear/p2p"
	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
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

	cmd := exec.Command("man", "cat")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return nil
}
