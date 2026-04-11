package ws

import (
	"log"
	"strings"

	"github.com/gabriel-logan/twitch-ai-bot/internal/ai"
)

func generateAIResponse(conversation []ai.RequestMessage, apiKey, model, fallbackModel string, twitchMaxLength int, whoExecuted string) (string, error) {
	response, err := ai.CallGroq(conversation, apiKey, model)
	if err != nil {
		log.Println(whoExecuted+": primary model error: ", err)

		responseFb, errFb := ai.CallGroq(conversation, apiKey, fallbackModel)
		if errFb != nil {
			log.Println(whoExecuted+": fallback error: ", errFb)
			return "", errFb
		}

		response = responseFb
	}

	response = strings.TrimSpace(response)

	responseRunes := []rune(response)
	if len(responseRunes) > twitchMaxLength {
		response = string(responseRunes[:twitchMaxLength])
	}

	return response, nil
}
