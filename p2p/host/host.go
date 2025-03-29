package host

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"

	"github.com/b-sharman/pear/p2p"
)

func Start(roomid string) error {
	u, _ := url.Parse(p2p.ServerUrl)
	u = u.JoinPath("register", roomid)
	resp, err := http.Post(u.String(), "", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("name server returned %s\n", resp.Status)
	}

	cmd := exec.Command("tmux", "new", "-s", roomid)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tmux error: %s", err.Error())
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill)

	<-c
	killCmd := exec.Command("tmux", "kill-session", "-t", roomid)

	if err := killCmd.Run(); err != nil {
		return fmt.Errorf("tmux kill error: %s", err.Error())
	}

	return nil
}
