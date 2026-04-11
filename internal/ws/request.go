package ws

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
)

func registerEventSub(sessionID, eventSubType string) {
	env := config.GetEnv()

	token := storage.GetOauthToken()

	body := EventSubRequest{
		Type:    eventSubType,
		Version: "1",
		Condition: EventSubRequestCondition{
			BroadcasterUserID: env.TwitchBroadcasterID,
			UserID:            env.TwitchBotUserID,
		},
		Transport: EventSubRequestTransport{
			Method:    "websocket",
			SessionID: sessionID,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Println("[registerEventSub] error marshaling request body: ", err)
		return
	}

	const baseURL = "https://api.twitch.tv/helix/eventsub/subscriptions"

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		log.Println("[registerEventSub] error creating HTTP request: ", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", env.TwitchClientID)
	req.Header.Set("Content-Type", "application/json")

	log.Println("[registerEventSub] registering subscription type: ", eventSubType, " sessionID: ", sessionID)

	resp, err := clientHttp.Do(req)
	if err != nil {
		log.Println("[registerEventSub] HTTP request failed: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		log.Println("[registerEventSub] Twitch API returned status: ", resp.Status)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("[registerEventSub] error reading response body: ", err)
		}

		log.Println("[registerEventSub] response body: ", string(bodyBytes))
		return
	}

	log.Println("[registerEventSub] subscription created successfully - type: ", eventSubType, " status: ", resp.Status)
}

func sendMessage(message string) {
	env := config.GetEnv()

	token := storage.GetOauthToken()

	body := SendMessageRequest{
		BroadcasterID: env.TwitchBroadcasterID,
		SenderID:      env.TwitchBotUserID,
		Message:       message,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Println("[sendMessage] error marshaling request body: ", err)
		return
	}

	const baseURL = "https://api.twitch.tv/helix/chat/messages"

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		log.Println("[sendMessage] error creating HTTP request: ", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", env.TwitchClientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := clientHttp.Do(req)
	if err != nil {
		log.Println("[sendMessage] HTTP request failed: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Println("[sendMessage] Twitch API returned status: ", resp.Status)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("[sendMessage] error reading response body: ", err)
		}

		log.Println("[sendMessage] response body: ", string(bodyBytes))
		return
	}
}
