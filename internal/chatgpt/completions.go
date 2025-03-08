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
const model = "gpt-3.5-turbo"

// Structs for API request and response
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string      `json:"model"`
	Messages    []Message   `json:"messages"`
	Max_tokens  int         `json:"max_tokens"`
	N           int         `json:"n"`
	Stop        interface{} `json:"stop"`
	Temperature float64     `json:"temperature"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Function to send request to OpenAI API
// func GetCompletions(apiKey, question string) (string, error) {
func GetCompletions(apiKey string, messages []Message) (string, error) {

	requestData := ChatRequest{
		Model:       model,
		Messages:    messages,
		Max_tokens:  1024,
		N:           1,
		Stop:        nil,
		Temperature: 1,
	}

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
	if resp.StatusCode == http.StatusUnauthorized ||
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
