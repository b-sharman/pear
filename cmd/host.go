/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"os/exec"

	"github.com/spf13/cobra"
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

	shellCmd := exec.Command("tmux", "attach", "-t", "pear")
	shellCmd.Stdin = conn
	shellCmd.Stdout = conn
	err = shellCmd.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
