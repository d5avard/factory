package main

import (
	"fmt"

	"github.com/d5avard/factory/internal"
	"github.com/d5avard/factory/internal/chatgpt"
)

func main() {
	var filename string
	var err error

	if filename, err = internal.GetConfigFilename(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	config, err := internal.LoadConfig(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Fetch and print available models
	models, err := chatgpt.GetModels(config.APIKey)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Available OpenAI Models:")
	for _, model := range models {
		fmt.Println("-", model)
	}
}
