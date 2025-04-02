package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/d5avard/factory/internal"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var host string
var port string
var templatePath string
var staticPath string

var rootCmd = &cobra.Command{
	Use:   "web",
	Short: "Run the web HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		router := mux.NewRouter()
		router.HandleFunc("/", homeHandler).Methods("GET")
		router.HandleFunc("/status", statusHandler).Methods("GET")
		http.Handle("/", router)

		log.Printf("Starting server on %s", port)

		addr := fmt.Sprintf("%s:%s", host, port)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Could not start server: %s", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&host, "host", "0.0.0.0", "Host to bind the server to")
	rootCmd.Flags().StringVar(&port, "port", "80", "Port to run the server on")
	rootCmd.Flags().StringVar(
		&templatePath,
		"templatePath",
		"./templates",
		"Path to the HTML template")
	rootCmd.Flags().StringVar(
		&staticPath,
		"staticPath",
		"./static",
		"Path to the static files")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	internal.LogRequest(r, http.StatusText(http.StatusOK), http.StatusOK)
}

const homeFile = "home.html"

func homeHandler(w http.ResponseWriter, r *http.Request) {
	file := filepath.Join(templatePath, homeFile)
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		internal.HttpError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, nil); err != nil {
		internal.HttpError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	internal.LogRequest(r, http.StatusText(http.StatusOK), http.StatusOK)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
