package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var port string

var rootCmd = &cobra.Command{
	Use:   "web",
	Short: "Run the web HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		router := mux.NewRouter()
		router.HandleFunc("/status", statusHandler).Methods("GET")
		http.Handle("/", router)

		log.Printf("Starting server on %s", port)

		addr := fmt.Sprintf("localhost:%s", port)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Could not start server: %s", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&port, "port", ":80", "Port to run the server on")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
