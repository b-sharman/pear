package host

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/b-sharman/pear/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func Start(roomid string, exitSignal chan int) error {
	// start a libp2p node that listens on a random local TCP port
	node, err := libp2p.New(
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
		libp2p.EnableAutoRelayWithStaticRelays([]peerstore.AddrInfo{{
			ID:    p2p.RelayPeerID,
			Addrs: p2p.RelayMultiAddrs(),
		}}),
	)
	if err != nil {
		panic(err)
	}

	// get opened port
	var port string
	for _, la := range node.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	err = registerRoom(port, node.ID(), roomid)
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
	fmt.Println("Created room " + roomid)
	<-exitSignal

	return nil
}

func registerRoom(port string, id peerstore.ID, roomid string) error {
	// retrieve public IP
	res, err := http.Get("https://checkip.amazonaws.com/")
	if err != nil {
		return err
	}
	ip, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// create addr string from public ip, port, and id
	b := bytes.NewBufferString(fmt.Sprintf("/dns4/pear-programming.fly.dev/tcp/%v/ip4/%s/tcp/%v/p2p/%s", 1337, strings.TrimSpace(string(ip)), port, id))

	u, _ := url.Parse(p2p.ServerUrl)
	u = u.JoinPath("register", roomid)
	resp, err := http.Post(u.String(), "", b)

	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("name server returned %s\n", resp.Status)
	}
	return nil
}
