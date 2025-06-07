package main

import (
	"errors"
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

	file, err := os.OpenFile("question.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open question.log file: %v", err)
	}
	defer file.Close()

	// Set output of logs to the file
	log.SetOutput(file)

	// Optional: Set log flags (timestamp, file, etc.)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("new question command started")

	// Ensure an argument is provided
	if len(os.Args) < 2 {
		log.Fatalln(errors.New("please provide a question as an argument"))
		return
	}

	question := strings.Join(os.Args[1:], " ")

	debug := internal.GetDebugVar()
	log.Println("debug mode:", debug)

	if filename, err = internal.GetConfigFilename(debug); err != nil {
		fmt.Println("error:", err)
		return
	}

	// Load config file
	config, err := internal.LoadConfig(filename)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("question:", question)

	var messages []chatgpt.Message

	messages = append(
		messages,
		chatgpt.Message{
			Role: "assistant",
			Content: `
			# IDENTITY
			You are an assistant.

			# GOAL
			From a question, you will write an answer.`,
		})

	messages = append(messages, chatgpt.Message{Role: "user", Content: string(question)})

	// Get response from ChatGPT
	attr := chatgpt.NewDefaultAttributes()
	attr.Max_completion_tokens = 1024
	attr.Temperature = 1
	response, err := chatgpt.GetCompletions(config.APIKey, messages, attr)
	if err != nil {
		log.Println("error:", err)
		return
	}

	// Print response
	fmt.Println("response:", response)
}
