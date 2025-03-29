package host

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

const serverUrl = "https://pear-programming.fly.dev"

func Start() error {
	// replace with randomly generated ID
	roomid := "oheaohea"

	u, _ := url.Parse(serverUrl)
	u = u.JoinPath("register", roomid)
	resp, err := http.Post(u.String(), "", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("name server returned %s\n", resp.Status)
	}

	cmd := exec.Command("tmux")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tmux error: %s", err.Error())
	}

	return nil
}
