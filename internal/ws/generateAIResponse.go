package ws

import (
	"log"
	"strings"

	"github.com/gabriel-logan/twitch-ai-bot/internal/ai"
)

type GenerateAIResponseArgs struct {
	conversation    []ai.RequestMessage
	apiKey          string
	model           string
	fallbackModel   string
	twitchMaxLength int
	whoExecuted     string
}

func generateAIResponse(args GenerateAIResponseArgs) (string, error) {
	response, err := ai.CallGroq(args.conversation, args.apiKey, args.model)
	if err != nil {
		log.Println(args.whoExecuted+": primary model error: ", err)

		responseFb, errFb := ai.CallGroq(args.conversation, args.apiKey, args.fallbackModel)
		if errFb != nil {
			log.Println(args.whoExecuted+": fallback error: ", errFb)
			return "", errFb
		}

		response = responseFb
	}

	response = strings.TrimSpace(response)

	responseRunes := []rune(response)
	if len(responseRunes) > args.twitchMaxLength {
		response = string(responseRunes[:args.twitchMaxLength])
	}

	return response, nil
}
