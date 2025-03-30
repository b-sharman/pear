/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Host a Pear session from your terminal",
	Run: host,
}

func init() {
	rootCmd.AddCommand(hostCmd)
}

func host(cmd *cobra.Command, args []string) {
	fmt.Println("host called")
}
