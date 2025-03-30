package host

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"os/exec"
	"strings"

	"github.com/b-sharman/pear/p2p"
	"github.com/creack/pty"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/term"
)

func Start(roomid string, exitSignal chan int) error {
	// start a libp2p node that listens on a random local TCP port

	relay, err := peer.AddrInfoFromP2pAddr(p2p.RelayMultiAddrs()[0])
	if err != nil {
		panic(err)
	}

	node, err := libp2p.New(
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
		libp2p.EnableAutoRelayWithStaticRelays([]peerstore.AddrInfo{*relay}),
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

	_, err = client.Reserve(context.Background(), node, *relay)
	if err != nil {
		fmt.Printf("Host failed to receive a relay reservation from relay.")
		panic(err)
	}

	node.SetStreamHandler("/connect/0.0.0", func(s network.Stream) {
		cmd := exec.Command("tmux", "attach-session", "-t", roomid)

		// Start the command with a pty.
		ptmx, err := pty.Start(cmd)
		if err != nil {
			panic(err)
		}
		// Make sure to close the pty at the end.
		defer func() { _ = ptmx.Close() }() // Best effort.

		// Copy stdin to the pty and the pty to stdout.
		// NOTE: The goroutine will keep reading until the next keystroke before returning.
		go func() { _, _ = io.Copy(ptmx, s) }()
		_, _ = io.Copy(s, ptmx)
	})

	fmt.Println("Created room " + roomid)
	fmt.Println("Starting server...")
	time.Sleep(time.Second * 4)
	cmd := exec.Command("tmux", "new-session", "-s", roomid)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	cmd.Run()
	term.Restore(int(os.Stdin.Fd()), oldState)

	<-exitSignal

	node.Close()

	if err := cleanUpRoomId(roomid); err != nil {
		panic(err)
	}

	return nil
}

func cleanUpRoomId(roomid string) error {
	u, _ := url.Parse(p2p.ServerUrl)
	u = u.JoinPath("register", roomid)
	_, err := http.NewRequest("DELETE", u.String(), nil)
	return err
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

	// create address string from public ip, port, and id
	b := bytes.NewBufferString(fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", strings.TrimSpace(string(ip)), port, id))

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
