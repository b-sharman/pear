package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// serverLaunchCmd represents the serverLaunch command
var serverLaunchCmd = &cobra.Command{
	Use:   "server-launch",
	Short: "Launch an instance of the central server",
	Run: func(cmd *cobra.Command, args []string) {
		const addr = ":8080"
		fmt.Println("listening on", addr)
		listenAndServe(addr)
	},
	// Hidden: true,
}

func init() {
	rootCmd.AddCommand(serverLaunchCmd)
}

func listenAndServe(addr string) {
	http.HandleFunc("GET /api/lookup", func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("roomid") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		roomid := r.URL.Query().Get("roomid")
		multiaddr := ""
		log.Printf("lookup: %s -  %s\n", roomid, multiaddr)
		w.WriteHeader(200)
		fmt.Fprint(w, multiaddr)
	})

        log.Fatalln(http.ListenAndServe(addr, nil))
}
