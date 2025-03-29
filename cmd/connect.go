package cmd

import (
	"github.com/b-sharman/pear/p2p/client"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a room",
	Args: cobra.ExactArgs(1),
	Run: connect,
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

func connect(cmd *cobra.Command, args []string) {
	client.Start(args[0])
}
