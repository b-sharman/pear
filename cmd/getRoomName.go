package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// getRoomNameCmd represents the getRoomName command
var getRoomNameCmd = &cobra.Command{
	Use:   "get-room-name",
	Short: "lets you see what the room name is that you are currently connected to",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("http://localhost:8923/roomid")
		if err != nil {
			fmt.Println("unable to get the roomid. likely no room is hosted: ", err.Error())
			return
		}
		io.Copy(os.Stdout, resp.Body)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(getRoomNameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gertRoomNameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gertRoomNameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
