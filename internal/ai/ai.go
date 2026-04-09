package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
)

type RequestMessage struct {
	Role    string `json:"role"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Request struct {
	Model    string           `json:"model"`
	Messages []RequestMessage `json:"messages"`
}

type ResponseChoiceMessage struct {
	Role    string  `json:"role"`
	Content *string `json:"content"`
}

type ResponseChoice struct {
	Message ResponseChoiceMessage `json:"message"`
}

type Response struct {
	Model   string           `json:"model"`
	Choices []ResponseChoice `json:"choices"`
}

func CallGroq(messages []RequestMessage) (string, error) {
	const url = "https://api.groq.com/openai/v1/chat/completions"

	env := config.GetEnv()

	apiKey := env.GroqAPIKey
	model := env.GroqModel

	return CallGenericAI(messages, apiKey, model, url)
}

func CallGenericAI(messages []RequestMessage, apiKey, model, url string) (string, error) {
	payload := Request{
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

	clientHttp := &http.Client{
		Timeout: config.GetEnv().ContextRequestDuration,
	}
	resp, err := clientHttp.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		log.Println("Result.Choices is empty")
		return "", nil
	}

	msg := result.Choices[0].Message.Content

	if msg == nil {
		log.Println("Result.Choices.Message.Content is nil")
		return "", nil
	}

	return *msg, nil
}
