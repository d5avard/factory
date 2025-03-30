package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/d5avard/factory/internal"
	"github.com/d5avard/factory/internal/chatgpt"
)

var config internal.Config
var messages []chatgpt.Message

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/chat.html")
	if err != nil {
		internal.HttpError(w, r, "Could not load template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		internal.HttpError(w, r, "Could not render template", http.StatusInternalServerError)
		return
	}
	internal.LogRequest(r, http.StatusText(http.StatusOK), http.StatusOK)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	question := r.URL.Query().Get("question")
	if question == "" {
		internal.HttpError(w, r, "No question provided", http.StatusBadRequest)
		return
	}

	messages = append(messages, chatgpt.Message{Role: "user", Content: question})
	attr := chatgpt.NewDefaultAttributes()
	attr.Max_completion_tokens = 1024
	attr.Temperature = 1
	answer, err := chatgpt.GetCompletions(config.APIKey, messages, attr)
	if err != nil {
		internal.HttpError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	messages = append(messages, chatgpt.Message{Role: "assistant", Content: answer})

	internal.LogRequest(r, http.StatusText(http.StatusOK), http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(answer))
}

func main() {
	debug := internal.GetDebugVar()
	log.Println("Debug mode:", debug)

	filename, err := internal.GetConfigFilename(debug)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	config, err = internal.LoadConfig(filename)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	messages = append(messages, chatgpt.Message{Role: "system", Content: "You are a helpful assistant."})

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/get", getHandler)
	port := ":8080"
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
