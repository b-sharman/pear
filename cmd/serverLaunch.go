package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// serverLaunchCmd represents the serverLaunch command
var serverLaunchCmd = &cobra.Command{
	Use:   "server-launch",
	Short: "Launch an instance of the central server",
	Run: func(cmd *cobra.Command, args []string) {
		const httpAddr = ":8080"
		const udpAddr = "0.0.0.0:4000"
		fmt.Println("listening on", httpAddr)
		go udpRelayServer(udpAddr)
		listenAndServeHTTP(httpAddr)
	},
	// Hidden: true,
}

func init() {
	rootCmd.AddCommand(serverLaunchCmd)
}

func listenAndServeHTTP(addr string) {
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
func udpRelayServer(addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	buffer := make([]byte, 65535)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error reading from UDP:", err)
			continue
		}
		if err == nil {
			go func() {
				for {
					_, err := conn.WriteTo([]byte("j"), remoteAddr)
					if err != nil {
						log.Println("Error writing to UDP:", err)
					}
					time.Sleep(time.Second * 4)
				}
			}()
		}

		// Write the received data directly to stdout
		_, _ = os.Stdout.Write(buffer[:n])
		if err != nil {
			log.Println("Error writing to stdout:", err)
		}

		// Optionally log the source address to stderr
		fmt.Fprintf(os.Stderr, "Received %d bytes from %s\n", n, remoteAddr)
	}
}
