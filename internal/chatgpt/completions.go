package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/d5avard/factory/internal"
)

const completionsURL = "https://api.openai.com/v1/chat/completions"

// Model:    "gpt-4", // You can use "gpt-3.5-turbo" for a cheaper option
const ModelGpt35Turbo = "gpt-3.5-turbo"
const ModelGpt40Mini = "gpt-4o-mini"

// Structs for API request and response
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Attributes struct {
	Model                 string      `json:"model"`
	Max_completion_tokens int         `json:"max_completion_tokens"`
	N                     int         `json:"n"`
	Stop                  interface{} `json:"stop"`
	Temperature           float64     `json:"temperature"`
}

func NewDefaultAttributes() *Attributes {
	return &Attributes{
		Model:                 ModelGpt40Mini,
		Max_completion_tokens: 1024,
		N:                     1,
		Stop:                  nil,
		Temperature:           1,
	}
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	// Max_tokens            int         `json:"max_tokens"`
	Max_completion_tokens int         `json:"max_completion_tokens"`
	N                     int         `json:"n"`
	Stop                  interface{} `json:"stop"`
	Temperature           float64     `json:"temperature"`
}

func NewChatRequest(messages []Message, attributes *Attributes) *ChatRequest {
	return &ChatRequest{
		Model:                 attributes.Model,
		Messages:              messages,
		Max_completion_tokens: attributes.Max_completion_tokens,
		N:                     attributes.N,
		Stop:                  attributes.Stop,
		Temperature:           attributes.Temperature,
	}
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Function to send request to OpenAI API
// func GetCompletions(apiKey, question string) (string, error) {
func GetCompletions(apiKey string, messages []Message, attributes *Attributes) (string, error) {

	requestData := NewChatRequest(messages, attributes)

	// Convert struct to JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// Make HTTP request
	req, err := http.NewRequest("POST", completionsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	internal.LogRequest(req, http.StatusText(resp.StatusCode), resp.StatusCode)

	// Check for errors
	// If an Api key is invalid
	// if the ChatGPT model is not available
	if resp.StatusCode == http.StatusBadRequest ||
		resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusNotFound {
		var error Error
		if err := json.Unmarshal(body, &error); err != nil {
			return "", err
		}
		return "", errors.New(error.Error.Message)
	}

	// Parse JSON response
	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return "", err
	}

	// Extract and return answer
	if len(chatResponse.Choices) > 0 {
		return chatResponse.Choices[0].Message.Content, nil
	}

	return "", errors.New("no response from chatgpt")
}
