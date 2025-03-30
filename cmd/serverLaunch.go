package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"slices"

	"github.com/b-sharman/pear/constants"

	"github.com/spf13/cobra"
)

// serverLaunchCmd represents the serverLaunch command
var serverLaunchCmd = &cobra.Command{
	Use:   "server-launch",
	Short: "Launch an instance of the central server",
	Run: func(cmd *cobra.Command, args []string) {
		const udpAddr = "0.0.0.0:4000"
		udpRelayServer(udpAddr)
	},
	// Hidden: true,
}

func init() {
	rootCmd.AddCommand(serverLaunchCmd)
}

func udpRelayServer(addr string) {
	log.Println("Starting server on", addr)

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	buffer := make([]byte, 65535)

	clients := make([]*net.UDPAddr, 0)
	var host *net.UDPAddr = nil

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error reading from UDP:", err)
			return
		}
		if remoteAddr == host {
			for _, client := range clients {
				conn.WriteToUDP(buffer[:n], client)
			}
		} else if (slices.Contains(clients, remoteAddr)) {
			conn.WriteToUDP(buffer[:n], host)
		} else {
			if n != 1 {
				log.Printf("Expected 1 byte but received %d bytes\n", n)
			}
			switch (buffer[0]) {
			case constants.CLIENT:
				clients = append(clients, remoteAddr)
				// TODO: need some way to tell the host to spawn another tmux session
			case constants.HOST:
				if host != nil {
					log.Println("Uh oh, multiple nodes claim to be the host. Giving favor to the most recent.")
				}
				host = remoteAddr
			default:
				log.Printf("Expected CLIENT (%d) or HOST (%d) but got %d\n", constants.CLIENT, constants.HOST, buffer[0])
			}
		}

		// Write the received data directly to stdout
		_, err = os.Stdout.Write(buffer[:n])
		if err != nil {
			log.Println("Error writing to stdout:", err)
		}

		// Optionally log the source address to stderr
		fmt.Fprintf(os.Stderr, "Received %d bytes from %s\n", n, remoteAddr)
	}
}
