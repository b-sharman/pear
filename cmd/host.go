/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/b-sharman/pear/p2p/host"
	"github.com/spf13/cobra"
	"github.com/wordgen/wordgen"
	"github.com/wordgen/wordlists"
	"golang.org/x/text/language"
)

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Start a pear session",
	Run: func(cmd *cobra.Command, args []string) {
		gen := wordgen.NewGenerator()
		gen.Words = wordlists.EffLarge
		gen.Count = 3
		gen.Casing = "lower"
		gen.Separator = "-"
		gen.Language = language.English

		result, err := gen.Generate()
		if err != nil {
			fmt.Println("roomid generation error: ", err)
			return
		}

		if err = host.Start(result); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(hostCmd)
}
