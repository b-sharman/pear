package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
)

// genKeyCmd represents the genKey command
var genKeyCmd = &cobra.Command{
	Use: "genKey",
	Run: func(cmd *cobra.Command, args []string) {
		err := func() error {
			priv, pub, err := crypto.GenerateEd25519Key(rand.Reader)
			if err != nil {
				return err
			}

			rawPrivKey, err := priv.Raw()
			if err != nil {
				return err
			}

			id, err := peer.IDFromPublicKey(pub)
			if err != nil {
				return err
			}

			fmt.Print("id: \n\n")
			fmt.Println(id)

			fmt.Print("raw private key: \n\n")
			fmt.Println(hex.EncodeToString(rawPrivKey))

			return nil
		}()
		if err != nil {
			log.Fatalln(err)
		}
	},
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(genKeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genKeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
