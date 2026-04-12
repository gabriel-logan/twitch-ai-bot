package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var clientHttp = &http.Client{
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 100,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		MaxConnsPerHost:     100,
	},
}

func CallGroq(messages []RequestMessage, apiKey, model string) (string, error) {
	const url = "https://api.groq.com/openai/v1/chat/completions"

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

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := clientHttp.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error when trying to read response body: ", err)
		}

		statusCodeString := strconv.Itoa(resp.StatusCode)

		return "", errors.New("API error: " + string(bodyBytes) + " - " + statusCodeString)
	}

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("No choices returned from API - empty")
	}

	msg := result.Choices[0].Message.Content

	if msg == "" {
		return "", errors.New("No content returned from API - empty")
	}

	return msg, nil
}
