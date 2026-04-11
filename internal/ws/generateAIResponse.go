package ws

import (
	"fmt"
	"log"
	"strings"

	"github.com/gabriel-logan/twitch-ai-bot/internal/ai"
)

type GenerateAIResponseArgs struct {
	conversation    []ai.RequestMessage
	apiKey          string
	models          []string
	twitchMaxLength int
	whoExecuted     string
}

func generateAIResponse(args GenerateAIResponseArgs) (string, error) {
	var lastErr error

	if len(args.models) == 0 {
		return "", fmt.Errorf("no models provided")
	}

	for i, model := range args.models {
		response, err := ai.CallGroq(args.conversation, args.apiKey, model)
		if err != nil {
			log.Printf("%s: model %d (%s) error: %v\n", args.whoExecuted, i, model, err)
			lastErr = err
			continue
		}

		response = strings.TrimSpace(response)

		responseRunes := []rune(response)
		if len(responseRunes) > args.twitchMaxLength {
			response = string(responseRunes[:args.twitchMaxLength])
		}

		return response, nil
	}

	return "", fmt.Errorf("all models failed, last error: %w", lastErr)
}
