package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/d5avard/factory/internal/chatgpt"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var port string

var rootCmd = &cobra.Command{
	Use:   "requester",
	Short: "Run the requester HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		router := mux.NewRouter()
		router.HandleFunc("/status", statusHandler).Methods("GET")
		// router.HandleFunc("/messages", messagesHandler).Methods("POST")
		http.Handle("/", router)
		log.Printf("Starting server on %s", port)
		// if err := http.ListenAndServe(port, nil); err != nil {
		addr := fmt.Sprintf("localhost:%s", port)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Could not start server: %s", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&port, "port", ":80", "Port to run the server on")
}

type Response struct {
	Message string `json:"message"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

// messagesHandler responds to GET requests at /messages.
func messagesHandler(w http.ResponseWriter, r *http.Request) {

	// Based on the structure Message in internal/chatgpt/completions.go.
	// Validate JSON
	// Validate it's a list of messages
	// Validate each message has a role and content
	// Validate each role is either "user" or "assistant"
	// Validate each content is a string
	// Validate each content is not empty
	// Validate each content is not too long
	// Validate each content is not too short
	// Define a Message type for validation.
	// type Message struct {
	// 	Role    string `json:"role"`
	// 	Content string `json:"content"`
	// }

	const (
		minContentLength = 1
		maxContentLength = 4096
	)

	// Decode JSON as a slice of Message.
	var messages []chatgpt.Message
	if err := json.NewDecoder(r.Body).Decode(&messages); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Ensure the JSON is a non-empty list.
	if len(messages) == 0 {
		http.Error(w, "Expected a non-empty list of messages", http.StatusBadRequest)
		return
	}

	for _, msg := range messages {
		// Validate that each message has a valid role.
		if msg.Role != "user" && msg.Role != "assistant" {
			http.Error(w, "Message role must be 'user' or 'assistant'", http.StatusBadRequest)
			return
		}
		// Validate that the content is a non-empty string.
		if msg.Content == "" {
			http.Error(w, "Message content cannot be empty", http.StatusBadRequest)
			return
		}
		// Validate content length.
		if len(msg.Content) < minContentLength {
			http.Error(w, "Message content is too short", http.StatusBadRequest)
			return
		}
		if len(msg.Content) > maxContentLength {
			http.Error(w, "Message content is too long", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "This is the /messages endpoint"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
