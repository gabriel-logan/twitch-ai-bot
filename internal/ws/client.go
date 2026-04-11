package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gabriel-logan/twitch-ai-bot/internal/ai"
	"github.com/gabriel-logan/twitch-ai-bot/internal/config"
	"github.com/gabriel-logan/twitch-ai-bot/internal/helper"
	"github.com/gabriel-logan/twitch-ai-bot/internal/storage"
	"github.com/gorilla/websocket"
)

var (
	clientHttp = &http.Client{
		Timeout: 12 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
	}
	ctx       context.Context
	ctxcancel context.CancelFunc
)

func StartBot() {
	ctx, ctxcancel = context.WithCancel(context.Background())

	for {
		select {
		case <-ctx.Done():
			return
		default:
			run(ctx)

			if ctx.Err() != nil {
				log.Println("Closing connection...: ", ctx.Err())
				return
			}

			log.Println("Reconnecting... (5 seconds)")

			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
		}
	}
}

func StopBot() {
	if ctxcancel != nil {
		ctxcancel()

		storage.SetBotIsOn(false)

		log.Println("Bot stopped")
	}
}

func run(ctx context.Context) {
	const twitchWS = "wss://eventsub.wss.twitch.tv/ws"

	conn, response, err := websocket.DefaultDialer.Dial(twitchWS, nil)
	if err != nil {
		log.Printf("Error when trying to connect to %s: %v \n", twitchWS, err)
		log.Printf("Response: %v \n", response)
		return
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	listenTwitch(ctx, conn)
}

func listenTwitch(ctx context.Context, conn *websocket.Conn) { // nosonar
	env := config.GetEnv()

	const twitchMaxLength = 500

	var sessionID string

	systemTxt, err := helper.LoadFile("system_prompt.txt")
	if err != nil {
		log.Println("✖ Error loading system_prompt.txt: ", err)
		return
	}

	initialSystemPrompt := systemTxt

	const defaultMsg = "Don't create very long messages; messages should be short, with a maximum of 480 characters. You were created by Gabriel Logan - https://github.com/gabriel-logan, in case someone asks a related question. Always reply in the same language as the user who is speaking."

	initialSystemPrompt = initialSystemPrompt + defaultMsg + "Your name is defined as " + env.TwitchKeyWordToCallBot

	if env.GroqMaxContextInput < 5 {
		log.Println("GroqMaxContextInput must be at least 5")
		return
	}

	conversation := NewConversation(env.GroqMaxContextInput, ai.RequestMessage{
		Role:    "system",
		Name:    env.AppName,
		Content: initialSystemPrompt,
	})

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			log.Println("read error: ", err)
			return
		}

		var data WSMessage
		if err := json.Unmarshal(msg, &data); err != nil {
			log.Println("json error: ", err)
			continue
		}

		switch data.Metadata.MessageType {

		case "session_welcome":
			sessionID = data.Payload.Session.ID

			log.Println("Session ID:", sessionID)

			eventSubTypes := []string{
				"channel.chat.message",
				"channel.chat.notification",
			}

			for _, eventSubType := range eventSubTypes {
				go registerEventSub(sessionID, eventSubType)
			}

			storage.SetBotIsOn(true)

		case "notification":
			if data.Payload.Event.ChatterUserLogin == env.TwitchBotUserName {
				continue
			}

			if data.Metadata.SubscriptionType == "channel.chat.message" {
				msg := strings.TrimSpace(data.Payload.Event.Message.Text)

				if strings.ToLower(msg) == "ping" {
					sendMessage("pong")
					continue
				}

				if strings.Contains(strings.ToLower(msg), env.TwitchKeyWordToCallBot) {
					user := data.Payload.Event.ChatterUserLogin

					conversation.Add(ai.RequestMessage{
						Role:    "user",
						Name:    user,
						Content: msg,
					})

					response, err := generateAIResponse(conversation.BuildMessages(), env.GroqAPIKey, env.GroqModel, env.GroqModelFallback, twitchMaxLength, "message")
					if err != nil {
						sendMessage("Something went wrong!!!")
						continue
					}

					conversation.Add(ai.RequestMessage{
						Role:    "assistant",
						Name:    env.TwitchKeyWordToCallBot,
						Content: response,
					})

					sendMessage(response)
					continue
				}
			}

			if data.Metadata.SubscriptionType == "channel.chat.notification" {
				conversation.Add(ai.RequestMessage{
					Role:    "user",
					Name:    "notification",
					Content: data.Payload.Event.SystemMessage + " Respond to the user based on this. More info if exists: " + data.Payload.Event.Message.Text,
				})

				response, err := generateAIResponse(conversation.BuildMessages(), env.GroqAPIKey, env.GroqModel, env.GroqModelFallback, twitchMaxLength, "notification")
				if err != nil {
					continue
				}

				conversation.Add(ai.RequestMessage{
					Role:    "assistant",
					Name:    env.TwitchKeyWordToCallBot,
					Content: response,
				})

				sendMessage(response)
				continue
			}
		}
	}
}

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
		log.Println("json error: ", err)
		return
	}

	const baseURL = "https://api.twitch.tv/helix/eventsub/subscriptions"

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		log.Println("eventsub error: ", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", env.TwitchClientID)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Registering eventsub type: %s \n", eventSubType)

	resp, err := clientHttp.Do(req)
	if err != nil {
		log.Println("eventsub error: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		log.Println("eventsub error: ", resp.Status)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("eventsub error: ", err)
		}

		log.Println("eventsub response: ", string(bodyBytes))
		return
	}

	log.Printf("EventSub '%s' - status: %v", eventSubType, resp.Status)
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
		log.Println("json error: ", err)
		return
	}

	const baseURL = "https://api.twitch.tv/helix/chat/messages"

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		log.Println("send message error: ", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", env.TwitchClientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := clientHttp.Do(req)
	if err != nil {
		log.Println("send message error: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Println("send message failed: ", resp.Status)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("error reading response body: ", err)
		}

		log.Println("response body: ", string(bodyBytes))

		return
	}
}
