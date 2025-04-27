package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/d5avard/factory/internal"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var host string
var httpPort string
var tlsPort string
var templatePath string
var staticPath string
var certFile string
var keyFile string

var rootCmd = &cobra.Command{
	Use:   "web",
	Short: "Run the web HTTP server",
	Run:   run,
}

func init() {
	rootCmd.PersistentFlags().String("host", "0.0.0.0", "Host to bind the server to")
	rootCmd.PersistentFlags().String("httpPort", "80", "Port to run the server on")
	rootCmd.PersistentFlags().String("tlsPort", "443", "Port to run the TLS server on")
	rootCmd.PersistentFlags().String("certFile", "./certs/server.crt", "Path to the TLS certificate file")
	rootCmd.PersistentFlags().String("keyFile", "./certs/server.key", "Path to the TLS key file")
	rootCmd.PersistentFlags().String("templatePath", "./templates", "Path to the HTML template")
	rootCmd.PersistentFlags().String("staticPath", "./static", "Path to the static files")

	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("httpPort", rootCmd.PersistentFlags().Lookup("httpPort"))
	viper.BindPFlag("tlsPort", rootCmd.PersistentFlags().Lookup("tlsPort"))
	viper.BindPFlag("certFile", rootCmd.PersistentFlags().Lookup("certFile"))
	viper.BindPFlag("keyFile", rootCmd.PersistentFlags().Lookup("keyFile"))
	viper.BindPFlag("templatePath", rootCmd.PersistentFlags().Lookup("templatePath"))
	viper.BindPFlag("staticPath", rootCmd.PersistentFlags().Lookup("staticPath"))

	viper.BindEnv("host", "WEB_HOST")
	viper.BindEnv("httpPort", "WEB_HTTP_PORT")
	viper.BindEnv("tlsPort", "WEB_TLS_PORT")
	viper.BindEnv("certFile", "WEB_CERT_FILE")
	viper.BindEnv("keyFile", "WEB_KEY_FILE")
	viper.BindEnv("templatePath", "WEB_TEMPLATE_PATH")
	viper.BindEnv("staticPath", "WEB_STATIC_PATH")
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

func run(cmd *cobra.Command, args []string) {

	host = viper.GetString("host")
	httpPort = viper.GetString("httpPort")
	tlsPort = viper.GetString("tlsPort")
	templatePath = viper.GetString("templatePath")
	staticPath = viper.GetString("staticPath")
	certFile = viper.GetString("certFile")
	keyFile = viper.GetString("keyFile")

	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/", homeHandler).Methods("GET")
		router.HandleFunc("/status", statusHandler).Methods("GET")

		log.Printf("Starting server on %s", tlsPort)
		log.Printf("Read Certificate file %s", certFile)
		log.Printf("Read Private key file %s", keyFile)
		log.Printf("Read Template path %s", templatePath)
		log.Printf("Read Static path %s", staticPath)
		addr := fmt.Sprintf("%s:%s", host, tlsPort)
		err := http.ListenAndServeTLS(addr, certFile, keyFile, router)
		if err != nil {
			log.Fatalf("HTTPS server error: %v", err)
		}
	}()
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var target string
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			target = fmt.Sprintf("https://%s%s", r.Host, r.URL.RequestURI())
		} else {
			target = fmt.Sprintf("https://%s:%s:%s", host, tlsPort, r.URL.RequestURI())
		}

		log.Printf("Redirecting to HTTPS %s", target)
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	})

	log.Printf("Starting server on %s", httpPort)
	addr := fmt.Sprintf("%s:%s", host, httpPort)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
