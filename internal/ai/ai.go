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
	json.NewDecoder(resp.Body).Decode(&result)

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("API response error: %v", result)
	}

	msg := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := msg["content"].(string)

	return content, nil
}

func CallGenericAI(messages []Message, url string) (string, error) {
	// TODO: implement
	return "", nil
}
