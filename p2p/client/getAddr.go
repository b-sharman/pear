package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/b-sharman/pear/p2p"
)

func getAddr(roomid string) (string, error) {
	u, _ := url.Parse(p2p.ServerUrl)
	u = u.JoinPath("lookup")
	q := u.Query()
	q.Set("roomid", roomid)
	u.RawQuery = q.Encode()
	// get the desired multiaddr as a string from the registry
	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("name server returned: %s\n", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(body), nil
}
