/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/b-sharman/pear/p2p/host"

	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Start a pear session",
	Run: func(cmd *cobra.Command, args []string) {
		err := host.Start()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(hostCmd)
}
