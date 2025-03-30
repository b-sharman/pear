package cmd

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)
var hostname string

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Host a Pear session from your terminal",
	Run:   host,
}

func init() {
	rootCmd.AddCommand(hostCmd)
	hostCmd.Flags().StringVarP(&hostname, "dial", "d", "pear-programming-dial.fly.dev:3000", "")
}

func host(cmd *cobra.Command, args []string) {
	fmt.Println("host called")
	conn, err := net.Dial("udp", hostname)
	if err != nil {
		log.Println(err.Error())
	}

	c := exec.Command("tmux", "attach", "-t", "pear")

        // Start the command with a pty.
        ptmx, err := pty.Start(c)
        if err != nil {
		panic(err)
        }
        // Make sure to close the pty at the end.
        defer func() { _ = ptmx.Close() }() // Best effort.

        // Handle pty size.
        ch := make(chan os.Signal, 1)
        signal.Notify(ch, syscall.SIGWINCH)
        go func() {
                for range ch {
                        if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
                                log.Printf("error resizing pty: %s", err)
                        }
                }
        }()
        ch <- syscall.SIGWINCH // Initial resize.
        defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

        // Set stdin in raw mode.
        oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
        if err != nil {
                panic(err)
        }
        defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

        // Copy stdin to the pty and the pty to stdout.
        // NOTE: The goroutine will keep reading until the next keystroke before returning.
	writer := io.MultiWriter(ptmx, conn)
        go func() { _, _ = io.Copy(writer, os.Stdin) }()
        go func() { _, _ = io.Copy(os.Stdout, ptmx) }()

        _, _ = io.Copy(os.Stdout, conn)
}
