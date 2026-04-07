package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Payload struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func CallGroq(messages []Message) (string, error) {
	const url = "https://api.groq.com/openai/v1/chat/completions"

	env := config.GetEnv()

	apiKey := env.GroqAPIKey
	model := env.GroqModel

	return CallGenericAI(messages, apiKey, model, url)
}

func CallGenericAI(messages []Message, apiKey, model, url string) (string, error) {
	payload := Payload{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("API response error: %v", result)
	}

	msg := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := msg["content"].(string)

	return content, nil
}
