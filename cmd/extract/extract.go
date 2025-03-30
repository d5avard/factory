package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/d5avard/factory/internal"
	"github.com/d5avard/factory/internal/chatgpt"
)

var config internal.Config
var messages []chatgpt.Message

// args
// debug
// with full text
var fulltext bool = true

func main() {
	debug := internal.GetDebugVar()
	log.Println("Debug mode:", debug)

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading input: %v\n", err)
	}

	filename, err := internal.GetConfigFilename(debug)
	if err != nil {
		log.Fatalf("Error reading config file: %v\n", err)
	}

	config, err = internal.LoadConfig(filename)
	if err != nil {
		log.Fatalf("Error reading config file: %v\n", err)
	}

	if len(input) > 0 {
		maxTokensEstimated := len(input) / 4
		log.Println("Max tokens estimated:", maxTokensEstimated)
	}

	messages = append(
		messages,
		chatgpt.Message{
			Role: "assistant",
			Content: `
			# IDENTITY
			You are an assistant.

			# GOAL
			From a text, you will extract the content.
			You will resume the text in 256 words.
			You will write the title, the authors, the source, the date and the resume.
			You will output the text in Markdown format.`,
		})

	messages = append(messages, chatgpt.Message{Role: "user", Content: string(input)})
	attr := chatgpt.NewDefaultAttributes()
	attr.Max_completion_tokens = 4048
	attr.Temperature = 1
	answer, err := chatgpt.GetCompletions(config.APIKey, messages, attr)
	if err != nil {
		log.Println(err)
		return
	}

	if len(answer) > 0 {
		maxTokensOutputEstimated := len(answer) / 4
		log.Println("Max tokens ouput estimated:", maxTokensOutputEstimated)
	}

	if fulltext {
		fmt.Println(answer + "\n\n## Full text\n\n" + string(input))
		return
	}

	// Output the result in markdown format.
	fmt.Println(answer)
}
