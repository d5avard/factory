package chatgpt

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// OpenAI API URL for listing models
const modelsURL = "https://api.openai.com/v1/models"

// Struct to parse the response
type ModelList struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

// Function to get available models
func GetModels(apiKey string) ([]string, error) {
	// Create HTTP request
	req, err := http.NewRequest("GET", modelsURL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		var error Error
		if err := json.Unmarshal(body, &error); err != nil {
			return nil, err
		}
		return nil, errors.New(error.Error.Message)
	}

	// Parse JSON response
	var modelList ModelList
	err = json.Unmarshal(body, &modelList)
	if err != nil {
		return nil, err
	}

	// Extract model names
	models := []string{}
	for _, model := range modelList.Data {
		models = append(models, model.ID)
	}

	return models, nil
}
