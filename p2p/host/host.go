package host

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
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
		libp2p.EnableAutoRelayWithStaticRelays([]peer.AddrInfo{*relay}),
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

	// multiaddr.String -> registered username
	type peerUsername struct {
		addr     string
		username string
	}

	usernameEvent := make(chan peerUsername)
	usernames := map[string]string{}
	node.SetStreamHandler("/username/0.0.0", func(s network.Stream) {
		reader := bufio.NewReader(s)

		peer := s.Conn().RemoteMultiaddr().Multiaddr().String()
		go func() {
			name, err := reader.ReadString('\n')
			if err != nil {
				panic("Unable to read from stream!")
			}

			usernames[peer] = name
			usernameEvent <- peerUsername{peer, name}
		}()
	})

	// multiaddr.String -> registered username
	type peerSize struct {
		addr string
		size pty.Winsize
	}
	resizeEvent := make(chan peerSize)
	node.SetStreamHandler("/resize/0.0.0", func(s network.Stream) {
		reader := bufio.NewReader(s)

		peer := s.Conn().RemoteMultiaddr().Multiaddr().String()
		size := pty.Winsize{}
		go func() {
			b, err := reader.ReadBytes('\n')
			if err != nil {
				panic("Unable to read from stream!")
			}

			fmt.Println(b)
			if err = json.Unmarshal(b, &size); err != nil {
				panic("Unable to unmarshal size")
			}

			fmt.Println("Resize requested!")
			resizeEvent <- peerSize{peer, size}
		}()
	})

	node.SetStreamHandler("/connect/0.0.0", func(s network.Stream) {
		cmd := exec.Command("tmux", "attach-session", "-t", roomid)

		// Start the command with a pty.
		ptmx, err := pty.Start(cmd)
		if err != nil {
			panic(err)
		}

		peer := s.Conn().RemoteMultiaddr().Multiaddr().String()
		go func() {
			for pSize := range resizeEvent {
				if pSize.addr == peer {
					fmt.Println("Resizing!")
					if err := pty.Setsize(ptmx, &pSize.size); err != nil {
						fmt.Printf("error resizing pty: %s", err)
					}
				}
			}
		}()

		// Read and write standard input and standard output from and to the stream
		go func() { _, _ = io.Copy(ptmx, s) }()
		_, _ = io.Copy(s, ptmx)

		// Make sure to close the pty at the end.
		defer ptmx.Close()
		defer s.Reset()
		defer delete(usernames, peer)
	})

	fmt.Println("Created room " + roomid)
	go func() {
		http.HandleFunc("GET /roomid", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(roomid)) })
		http.ListenAndServe(":8923", nil)
	}()
	go func() {
		for pName := range usernameEvent {
			fmt.Printf("%s connected with username %s\n", pName.addr, pName.username)
		}
	}()
	fmt.Println("Starting server...")
	time.Sleep(time.Second * 3)
	cmd := exec.Command("tmux", "new-session", "-s", roomid)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	cmd.Run()
	term.Restore(int(os.Stdin.Fd()), oldState)

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

func registerRoom(port string, id peer.ID, roomid string) error {
	// retrieve public IP
	res, err := http.Get("https://checkip.amazonaws.com/")
	if err != nil {
		return err
	}
	ip, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Create address string from public ip, port, and id
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
