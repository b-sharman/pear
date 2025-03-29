package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// nameServerCmd represents the nameServer command
var nameServerCmd = &cobra.Command{
	Use:  "name-server",
	Long: `a name server that makes it possible to match roomIDs to libp2p multiaddrs`,
	Run: func(cmd *cobra.Command, args []string) {

		kv, err := bolt.Open("./db/names.db", 0600, nil)
		if err != nil {
			panic(err.Error())
		}
		kv.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("roomids-addrs"))
			return err
		})
		addr := ":8080"
		fmt.Println("listening on ", addr)
		listenAndServe(addr, kv)
	},
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(nameServerCmd)
}

func listenAndServe(addr string, db *bolt.DB) {
	http.HandleFunc("GET /api/lookup", func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("roomid") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		roomid := r.URL.Query().Get("roomid")
		multiaddr := ""
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("roomids-addrs"))
			addr := b.Get([]byte(roomid))
			multiaddr = string(addr)
			return nil
		})
		if multiaddr == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(200)
		fmt.Fprint(w, multiaddr)
	})

	http.HandleFunc("POST /api/register/{roomid}", func(w http.ResponseWriter, r *http.Request) {
		roomid := r.PathValue("roomid")
		if roomid == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		multiaddr, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("roomids-addrs"))
			err := b.Put([]byte(roomid), multiaddr)
			return err
		})

		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("DELETE /api/register/{roomid}", func(w http.ResponseWriter, r *http.Request) {
		roomid := r.PathValue("roomid")
		if roomid == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("roomids-addrs"))
			err := b.Delete([]byte(roomid))
			return err
		})

		w.WriteHeader(http.StatusOK)
	})

	log.Fatalln(http.ListenAndServe(addr, nil))
}
