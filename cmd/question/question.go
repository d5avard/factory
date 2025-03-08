package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/d5avard/factory/internal"
	"github.com/d5avard/factory/internal/chatgpt"
)

func main() {
	var filename string
	var err error

	// Ensure an argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run chatgpt-cli.go \"Your question here\"")
		return
	}

	// Join arguments into a single question string
	question := strings.Join(os.Args[1:], " ")

	fmt.Println("Question:", question)

	debug := internal.GetDebugVar()
	log.Println("Debug mode:", debug)

	if filename, err = internal.GetConfigFilename(debug); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Load config file
	config, err := internal.LoadConfig(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Get response from ChatGPT
	response, err := chatgpt.GetCompletions(config.APIKey, []chatgpt.Message{})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print response
	fmt.Println("\nChatGPT Response:\n" + response)
}
