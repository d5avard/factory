package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/d5avard/factory/internal"
	"github.com/d5avard/factory/internal/chatgpt"
)

type Response struct {
	Message string `json:"message"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./web/index.html")
	if err != nil {
		http.Error(w, "Could not load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	question := r.URL.Query().Get("question")
	if question == "" {
		question = "No question provided"
	}
	log.Printf("Question: %s", question)

	answer, err := chatgpt.GetCompletions(config.APIKey, question)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(answer))
}

var config internal.Config

func main() {
	debug := internal.GetDebugVar()
	log.Println("Debug mode:", debug)

	filename, err := internal.GetConfigFilename(debug)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	config, err = internal.LoadConfig(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/get", getHandler)
	port := ":8080"
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
