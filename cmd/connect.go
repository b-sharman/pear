package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/b-sharman/pear/p2p/client"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a room",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			<-ch
			cancel()
		}()

		client.Start(ctx, args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
