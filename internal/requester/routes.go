package requester

import (
	"encoding/json"
	"net/http"

	"github.com/d5avard/factory/internal"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type StatusRoutes struct{}

func (r StatusRoutes) Register(router *mux.Router, logger *internal.Logger) {
	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received status request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"status": "OK"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Failed to encode response", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}).Methods("GET")
}

type MessagesRoutes struct{}

func (r MessagesRoutes) Register(router *mux.Router, logger *internal.Logger) {

	router.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received messages request", zap.String("method", r.Method), zap.String("path", r.URL.Path))
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"status": "OK"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Failed to encode response", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}).Methods("GET")

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

	// const (
	// 	minContentLength = 1
	// 	maxContentLength = 4096
	// )

	// Decode JSON as a slice of Message.
	// var messages []chatgpt.Message
	// if err := json.NewDecoder(r.Body).Decode(&messages); err != nil {
	// 	http.Error(w, "Invalid JSON", http.StatusBadRequest)
	// 	return
	// }

	// Ensure the JSON is a non-empty list.
	// if len(messages) == 0 {
	// 	http.Error(w, "Expected a non-empty list of messages", http.StatusBadRequest)
	// 	return
	// }

	// for _, msg := range messages {
	// 	// Validate that each message has a valid role.
	// 	if msg.Role != "user" && msg.Role != "assistant" {
	// 		http.Error(w, "Message role must be 'user' or 'assistant'", http.StatusBadRequest)
	// 		return
	// 	}
	// 	// Validate that the content is a non-empty string.
	// 	if msg.Content == "" {
	// 		http.Error(w, "Message content cannot be empty", http.StatusBadRequest)
	// 		return
	// 	}
	// 	// Validate content length.
	// 	if len(msg.Content) < minContentLength {
	// 		http.Error(w, "Message content is too short", http.StatusBadRequest)
	// 		return
	// 	}
	// 	if len(msg.Content) > maxContentLength {
	// 		http.Error(w, "Message content is too long", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	// w.Header().Set("Content-Type", "application/json")
	// response := Response{Message: "This is the /messages endpoint"}
	// json.NewEncoder(w).Encode(response)
}
